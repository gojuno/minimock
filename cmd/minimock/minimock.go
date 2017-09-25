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
	programOptions struct {
		Interfaces      []interfaceInfo
		Suffix          string
		OutputFile      string
		StructName      string
		ImportWithTests bool
	}

	generateOptions struct {
		InterfaceName      string
		PackageName        string
		OutputFileName     string
		StructName         string
		SourcePackage      string
		DestinationPackage string
	}

	interfaceInfo struct {
		Package string
		Name    string
	}

	visitor struct {
		gen             *generator.Generator
		methods         map[string]*types.Signature
		sourceInterface string
	}
)

func main() {
	opts := processFlags()

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

	for _, i := range opts.Interfaces {
		if opts.ImportWithTests {
			cfg.ImportWithTests(i.Package)
		} else {
			cfg.Import(i.Package)
		}
	}

	var outPackageRealPath string

	stat, err := os.Stat(opts.OutputFile)
	if err != nil {
		if !os.IsNotExist(err) {
			die("failed to get file info for %s: %v", opts.OutputFile, err)
		}

		if outPackageRealPath, err = generator.PackageAbsPath(opts.OutputFile); err != nil {
			die("failed to get real path for the %s: %v", opts.OutputFile, err)
		}
	} else {
		if stat.IsDir() {
			outPackageRealPath = opts.OutputFile
		} else {
			outPackageRealPath = filepath.Dir(opts.OutputFile)
		}
	}

	destImportPath, err := generator.PackageOf(outPackageRealPath)
	if err != nil {
		die("failed to detect import path of the %s: %v", outPackageRealPath, err)
	}

	if opts.ImportWithTests {
		cfg.ImportWithTests(destImportPath)
	} else {
		cfg.Import(destImportPath)
	}

	prog, err := cfg.Load()
	if err != nil {
		die("failed to load source code: %v", err)
	}

	if len(opts.Interfaces) == 1 && strings.HasSuffix(opts.OutputFile, ".go") { //legacy mode
		genOpts := generateOptions{
			SourcePackage:      opts.Interfaces[0].Package,
			DestinationPackage: destImportPath,
			InterfaceName:      opts.Interfaces[0].Name,
			StructName:         opts.StructName,
			OutputFileName:     opts.OutputFile,
			PackageName:        prog.Package(destImportPath).Pkg.Name(),
		}

		if err := generate(prog, genOpts); err != nil {
			die("failed to generate %s: %v", opts.OutputFile, err)
		}
	} else {
		for _, i := range opts.Interfaces {
			genOpts := generateOptions{
				PackageName:        prog.Package(destImportPath).Pkg.Name(),
				InterfaceName:      i.Name,
				StructName:         i.Name + "Mock",
				SourcePackage:      i.Package,
				DestinationPackage: destImportPath,
				OutputFileName:     filepath.Join(outPackageRealPath, i.Name+opts.Suffix),
			}

			if err := generate(prog, genOpts); err != nil {
				die("failed to generate %s: %v", genOpts.OutputFileName, err)
			}

			fmt.Printf("Generated file: %s\n", genOpts.OutputFileName)
		}
	}
}

func generate(prog *loader.Program, opts generateOptions) error {
	if err := os.Remove(opts.OutputFileName); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove output file %s: %v", opts.OutputFileName, err)
	}

	gen := generator.New(prog)
	gen.ImportWithAlias(opts.DestinationPackage, "")
	gen.SetPackageName(opts.PackageName)
	gen.SetVar("structName", opts.StructName)
	gen.SetVar("interfaceName", opts.InterfaceName)
	gen.SetVar("packagePath", opts.SourcePackage)
	gen.SetHeader(fmt.Sprintf(`DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.7
The original interface %q can be found in %s`, opts.InterfaceName, opts.SourcePackage))
	gen.SetDefaultParamsPrefix("p")
	gen.SetDefaultResultsPrefix("r")

	v := &visitor{
		gen:             gen,
		sourceInterface: opts.InterfaceName,
	}

	pkg := prog.Package(opts.SourcePackage)
	if pkg == nil {
		return fmt.Errorf("unable to load package: %s", opts.SourcePackage)
	}

	for _, file := range pkg.Files {
		ast.Walk(v, file)
	}

	if v.methods == nil {
		return fmt.Errorf("interface %s was not found in %s", opts.InterfaceName, opts.SourcePackage)
	}

	if len(v.methods) == 0 {
		return fmt.Errorf("empty interface: %s", opts.InterfaceName)
	}

	if err := gen.ProcessTemplate("interface", template, v.methods); err != nil {
		return err
	}

	if err := gen.WriteToFilename(opts.OutputFileName); err != nil {
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

		//{{$methodName}}MinimockCounter returns a count of {{$interfaceName}}.{{$methodName}} invocations
		func (m *{{$structName}}) {{$methodName}}MinimockCounter() uint64 {
			return atomic.LoadUint64(&m.{{$methodName}}Counter)
		}
	{{ end }}

	//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
	//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) ValidateCallCounters() {
		{{ range $methodName, $method := . }}
			if m.{{$methodName}}Func != nil && atomic.LoadUint64(&m.{{$methodName}}Counter) == 0 {
				m.t.Fatal("Expected call to {{$structName}}.{{$methodName}}")
			}
		{{ end }}
	}

	//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
	//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) CheckMocksCalled() {
		m.Finish()
	}

	//Finish checks that all mocked methods of the interface have been called at least once
	//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
	func (m *{{$structName}}) Finish() {
		m.MinimockFinish()
	}

	//MinimockFinish checks that all mocked methods of the interface have been called at least once
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

func processFlags() *programOptions {
	var (
		input       = flag.String("f", "", "DEPRECATED: input file or import path of the package that contains interface declaration")
		help        = flag.Bool("h", false, "show this help message")
		interfaces  = flag.String("i", "", "comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader")
		output      = flag.String("o", "", "destination file name to place the generated mock or path to destination package when multiple interfaces are given")
		packageName = flag.String("p", "", "DEPRECATED: destination package name")
		suffix      = flag.String("s", "_mock_test.go", "output file name suffix which is added to file names when multiple interfaces are given")
		sname       = flag.String("t", "", "DEPRECATED: mock struct name (default <interface name>Mock)")
		withTests   = flag.Bool("withTests", false, "parse *_test.go files in the source package")
	)

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *input != "" {
		fmt.Printf("DEPRECATED FLAG: -f\n")
	}

	if *packageName != "" {
		fmt.Printf("DEPRECATED FLAG: -p\n")
	}

	if *sname != "" {
		fmt.Printf("DEPRECATED FLAG: -t\n")
	}

	if *input != "" && *interfaces != "" {
		if *sname == "" {
			*sname = *interfaces + "Mock"
		}

		*interfaces = getImportPath(*input) + "." + *interfaces
	}

	if *interfaces == "" {
		die("missing required parameter: -i, use -h flag for help")
	}

	if *output == "" {
		die("missing required parameter: -o, use -h flag for help")
	}

	interfacesList := []interfaceInfo{}
	for _, i := range strings.Split(*interfaces, ",") {
		chunks := strings.Split(i, ".")
		if len(chunks) < 2 {
			die("invalid interface name: %s\nname should be in the form <import path>.<interface type>, i.e. io.Reader\n", i)
		}

		importPath := getImportPath(strings.Join(chunks[0:len(chunks)-1], "."))

		interfacesList = append(interfacesList, interfaceInfo{Package: importPath, Name: chunks[len(chunks)-1]})
	}

	return &programOptions{
		Interfaces:      interfacesList,
		OutputFile:      *output,
		StructName:      *sname,
		ImportWithTests: *withTests,
		Suffix:          *suffix,
	}
}

func getImportPath(realPath string) string {
	_, err := os.Stat(realPath)
	if err == nil {
		importPath, err := generator.PackageOf(realPath)
		if err != nil {
			die("failed to detect import path of the %s: %v", realPath, err)
		}

		return importPath
	}

	return realPath
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
