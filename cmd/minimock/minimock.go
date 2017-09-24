package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/gojuno/generator"
)

type (
	options struct {
		InputFile       string
		OutputFile      string
		InterfaceName   string
		StructName      string
		Package         string
		ImportWithTests bool
		All             bool
		AllSuffix       string
	}

	visitor struct {
		gen             *generator.Generator
		methods         map[string]*types.Signature
		sourceInterface string
	}

	allFinder struct {
		pkg             *loader.PackageInfo
		foundInterfaces []string
	}
)

func main() {
	opts := processFlags()
	sourcePackagePath := opts.InputFile

	_, err := os.Stat(sourcePackagePath)
	if err == nil {
		if sourcePackagePath, err = generator.PackageOf(sourcePackagePath); err != nil {
			die("failed to detect import path of the %s: %v", sourcePackagePath, err)
		}
	}

	var destPackagePath string
	if opts.All {
		if opts.OutputFile == "" {
			destPackagePath = sourcePackagePath
		}
	} else {
		outPackage := filepath.Dir(opts.OutputFile)
		if destPackagePath, err = generator.PackageOf(outPackage); err != nil {
			die("failed to detect import path of the %s: %v", outPackage, err)
		}
	}

	cfg := loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(string) bool { return false },
		TypeChecker: types.Config{
			IgnoreFuncBodies:         true,
			FakeImportC:              true,
			DisableUnusedImportCheck: true,
			Error: func(err error) {},
		},
	}

	if opts.ImportWithTests {
		cfg.ImportWithTests(sourcePackagePath)
	} else {
		cfg.Import(sourcePackagePath)
	}

	if destPackagePath != sourcePackagePath {
		cfg.Import(destPackagePath)
	}

	prog, err := cfg.Load()
	if err != nil {
		die("failed to load source code: %v", err)
	}

	pkg := prog.Package(sourcePackagePath)
	if opts.Package == "" {
		opts.Package = pkg.Pkg.Name()
	}

	if opts.All {
		finder := &allFinder{pkg: pkg}
		for _, file := range pkg.Files {
			ast.Walk(finder, file)
		}

		if len(finder.foundInterfaces) == 0 {
			die("no interfaces found in %s", sourcePackagePath)
		}

		destPath, err := generator.PackageAbsPath(destPackagePath)
		if err != nil {
			die("failed to get absolute path for the %s", destPackagePath)
		}

		currentDir, err := os.Getwd()
		if err != nil {
			die("failed to get current directory: %v", err)
		}

		for _, interfaceName := range finder.foundInterfaces {
			outputFileName := interfaceName + opts.AllSuffix
			o := options{
				Package:       opts.Package,
				InterfaceName: interfaceName,
				StructName:    interfaceName + "Mock",
				OutputFile:    filepath.Join(destPath, outputFileName),
			}
			if err := generate(sourcePackagePath, destPackagePath, o, prog); err != nil {
				die("failed to generate %s: %v", o.OutputFile, err)
			}

			printName, err := filepath.Rel(currentDir, o.OutputFile)
			if err != nil && !strings.HasPrefix(o.OutputFile, currentDir) {
				printName = o.OutputFile
			}

			fmt.Printf("Generated file: %s\n", printName)
		}
	} else {
		if err := generate(sourcePackagePath, destPackagePath, *opts, prog); err != nil {
			die("failed to generate %s: %v", opts.OutputFile, err)
		}
	}
}

func generate(sourcePackagePath, destPackagePath string, opts options, prog *loader.Program) error {
	if err := os.Remove(opts.OutputFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove output file: %v", err)
	}

	gen := generator.New(prog)
	gen.ImportWithAlias(destPackagePath, "")
	gen.SetPackageName(opts.Package)
	gen.SetVar("structName", opts.StructName)
	gen.SetVar("interfaceName", opts.InterfaceName)
	gen.SetVar("packagePath", sourcePackagePath)
	gen.SetHeader(fmt.Sprintf(`DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.7
The original interface %q can be found in %s`, opts.InterfaceName, sourcePackagePath))
	gen.SetDefaultParamsPrefix("p")
	gen.SetDefaultResultsPrefix("r")

	v := &visitor{
		gen:             gen,
		sourceInterface: opts.InterfaceName,
	}

	pkg := prog.Package(sourcePackagePath)
	if pkg == nil {
		return fmt.Errorf("unable to load package: %s", sourcePackagePath)
	}

	for _, file := range pkg.Files {
		ast.Walk(v, file)
	}

	if v.methods == nil {
		return fmt.Errorf("interface %s was not found in %s", opts.InterfaceName, sourcePackagePath)
	}

	if len(v.methods) == 0 {
		return fmt.Errorf("empty interface: %s", opts.InterfaceName)
	}

	if err := gen.ProcessTemplate("interface", template, v.methods); err != nil {
		return err
	}

	if err := gen.WriteToFilename(opts.OutputFile); err != nil {
		return err
	}

	return nil
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch ts := node.(type) {
	case *ast.FuncDecl:
		return nil
	case *ast.TypeSpec:
		exprType, err := v.gen.ExpressionType(ts.Type)
		if err != nil {
			die("failed to get expression for %T %s: %v", ts.Type, ts.Name.Name, err)
		}

		var i *types.Interface

		switch t := exprType.(type) {
		case *types.Named:
			underlying, ok := t.Underlying().(*types.Interface)
			if !ok {
				return nil
			}
			i = underlying
		case *types.Interface:
			i = t
		}

		if ts.Name.Name == v.sourceInterface {
			v.processInterface(i)
		}

		return nil
	}

	return v
}

func (v *visitor) processInterface(t *types.Interface) {
	v.methods = make(map[string]*types.Signature)

	for i := 0; i < t.NumMethods(); i++ {
		v.methods[t.Method(i).Name()] = t.Method(i).Type().(*types.Signature)
	}
}

func (af *allFinder) Visit(node ast.Node) ast.Visitor {
	switch ts := node.(type) {
	case *ast.FuncDecl:
		return nil
	case *ast.TypeSpec:
		exprType := af.pkg.TypeOf(ts.Type)
		if exprType == nil {
			die("failed to get expression type for %T %s", ts.Type, ts.Name.Name)
		}

		var i *types.Interface

		switch t := exprType.(type) {
		case *types.Named:
			underlying, ok := t.Underlying().(*types.Interface)
			if !ok {
				return nil
			}
			i = underlying
		case *types.Interface:
			i = t
		default:
			return nil
		}

		if !i.Empty() {
			af.foundInterfaces = append(af.foundInterfaces, ts.Name.Name)
		}

		return nil
	}

	return af
}

const template = `
	import "github.com/gojuno/minimock"

	//{{$structName}} implements {{$packagePath}}.{{$interfaceName}}
	type {{$structName}} struct {
		t minimock.Tester

		{{ range $methodName, $method := . }}
			{{$methodName}}Func func{{ signature $method }}
			{{$methodName}}Counter uint64
			{{$methodName}}Mock m{{$structName}}{{$methodName}}
		{{ end }}
	}

	//New{{$structName}} returns a mock for {{$packagePath}}.{{$interfaceName}}
	func New{{$structName}}(t minimock.Tester) *{{$structName}} {
		m := &{{$structName}}{t: t}

		if controller, ok := t.(minimock.MockController); ok {
			controller.RegisterMocker(m)
		}

		{{ range $methodName, $method := . }}m.{{$methodName}}Mock = m{{$structName}}{{$methodName}}{mock: m}
		{{ end }}

		return m
	}

	{{ range $methodName, $method := . }}
		type m{{$structName}}{{$methodName}} struct {
			mock *{{$structName}}
		}

		//Return sets up a mock for {{$interfaceName}}.{{$methodName}} to return Return's arguments
		func (m m{{$structName}}{{$methodName}}) Return({{results $method}}) *{{$structName}} {
			m.mock.{{$methodName}}Func = func({{params $method}}) ({{(results $method).Types}}) {
				return {{ (results $method).Names }}
			}
			return m.mock
		}

		//Set uses given function f as a mock of {{$interfaceName}}.{{$methodName}} method
		func (m m{{$structName}}{{$methodName}}) Set(f func({{params $method}}) ({{results $method}})) *{{$structName}}{
			m.mock.{{$methodName}}Func = f
			return m.mock
		}

		//{{$methodName}} implements {{$packagePath}}.{{$interfaceName}} interface
		func (m *{{$structName}}) {{$methodName}}{{signature $method}} {
			defer atomic.AddUint64(&m.{{$methodName}}Counter, 1)

			if m.{{$methodName}}Func == nil {
				m.t.Fatal("Unexpected call to {{$structName}}.{{$methodName}}")
				return
			}

			{{if gt (len (results $method)) 0 }}
			return {{ end }} m.{{$methodName}}Func({{(params $method).Pass}})
		}

		//Get{{$methodName}}Counter returns a count of {{$interfaceName}}.{{$methodName}} invocations
		func (m *{{$structName}}) Get{{$methodName}}Counter() uint64 {
			return atomic.LoadUint64(&m.{{$methodName}}Counter)
		}
	{{ end }}

	//ValidateCallCounters checks that all mocked methods of the iterface have been called at least once
	//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) ValidateCallCounters() {
		{{ range $methodName, $method := . }}
			if m.{{$methodName}}Func != nil && atomic.LoadUint64(&m.{{$methodName}}Counter) == 0 {
				m.t.Fatal("Expected call to {{$structName}}.{{$methodName}}")
			}
		{{ end }}
	}

	//CheckMocksCalled checks that all mocked methods of the iterface have been called at least once
	//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) CheckMocksCalled() {
		m.Finish()
	}

	//Finish checks that all mocked methods of the iterface have been called at least once
	//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
	func (m *{{$structName}}) Finish() {
		m.MinimockFinish()
	}

	//MinimockFinish checks that all mocked methods of the iterface have been called at least once
	func (m *{{$structName}}) MinimockFinish() {
		{{ range $methodName, $method := . }}
			if m.{{$methodName}}Func != nil && atomic.LoadUint64(&m.{{$methodName}}Counter) == 0 {
				m.t.Fatal("Expected call to {{$structName}}.{{$methodName}}")
			}
		{{ end }}
	}

	//Wait waits for all mocked methods to be called at least once
	//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
	func (m *{{$structName}}) Wait(timeout time.Duration) {
		m.MinimockWait(timeout)
	}

	//MinimockWait waits for all mocked methods to be called at least once
	//this method is called by minimock.Controller
	func (m *{{$structName}}) MinimockWait(timeout time.Duration) {
		timeoutCh := time.After(timeout)
		for {
			ok := true
			{{ range $methodName, $method := . }}ok = ok && (m.{{$methodName}}Func == nil || atomic.LoadUint64(&m.{{$methodName}}Counter) > 0)
			{{ end }}

			if ok {
				return
			}

			select {
			case <-timeoutCh:
				{{ range $methodName, $method := . }}
					if m.{{$methodName}}Func != nil && atomic.LoadUint64(&m.{{$methodName}}Counter) == 0 {
						m.t.Error("Expected call to {{$structName}}.{{$methodName}}")
					}
				{{ end }}
				m.t.Fatalf("Some mocks were not called on time: %s", timeout)
				return
			default:
				time.Sleep(time.Millisecond)
			}
		}
	}

	//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
	//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
	func (m *{{$structName}}) AllMocksCalled() bool {
		{{ range $methodName, $method := . }}
			if m.{{$methodName}}Func != nil && atomic.LoadUint64(&m.{{$methodName}}Counter) == 0 {
				return false
			}
		{{ end }}

		return true
	}`

func processFlags() *options {
	var (
		all       = flag.Bool("a", false, "generate mocks for all interfaces found in the file/package")
		allSuffix = flag.String("allSuffix", "_mock_test.go", "output file name suffix, ignored when -a flag is not set")
		input     = flag.String("f", "", "input file or import path of the package that contains interface declaration")
		name      = flag.String("i", "", "name of the interface to mock (ignored when -a flag is set)")
		output    = flag.String("o", "", "destination file name to place the generated mock or path to destination package when -a flag is set (source package is used by default)")
		pkg       = flag.String("p", "", "destination package name (source package name is used by default)")
		sname     = flag.String("t", "", "mock struct name (default <interface name>Mock)")
		withTests = flag.Bool("withTests", false, "parse *_test.go files in the source package")
	)

	flag.Parse()

	if *input == "" {
		fmt.Printf("missing required parameter -f")
		flag.Usage()
		os.Exit(1)
	}

	if !*all && (*name == "" || *output == "" || !strings.HasSuffix(*output, ".go")) {
		flag.Usage()
		os.Exit(1)
	}

	if *sname == "" {
		*sname = *name + "Mock"
	}

	return &options{
		InputFile:       *input,
		OutputFile:      *output,
		InterfaceName:   *name,
		Package:         *pkg,
		StructName:      *sname,
		ImportWithTests: *withTests,
		All:             *all,
		AllSuffix:       *allSuffix,
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
