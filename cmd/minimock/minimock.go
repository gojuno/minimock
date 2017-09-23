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
	}

	visitor struct {
		gen             *generator.Generator
		methods         map[string]*types.Signature
		sourceInterface string
	}
)

func main() {
	opts := processFlags()
	packagePath := opts.InputFile

	if _, err := os.Stat(packagePath); err == nil {
		if packagePath, err = generator.PackageOf(packagePath); err != nil {
			die(err)
		}
	}

	destPackagePath, err := generator.PackageOf(filepath.Dir(opts.OutputFile))
	if err != nil {
		die(err)
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
		cfg.ImportWithTests(packagePath)
	} else {
		cfg.Import(packagePath)
	}

	if err := os.Remove(opts.OutputFile); err != nil && !os.IsNotExist(err) {
		die(err)
	}

	if destPackagePath != packagePath {
		cfg.Import(destPackagePath)
	}

	prog, err := cfg.Load()
	if err != nil {
		die(err)
	}

	gen := generator.New(prog)
	gen.ImportWithAlias(destPackagePath, "")
	gen.SetPackageName(opts.Package)
	gen.SetVar("structName", opts.StructName)
	gen.SetVar("interfaceName", opts.InterfaceName)
	gen.SetVar("packagePath", packagePath)
	gen.SetHeader(fmt.Sprintf(`DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.5
The original interface %q can be found in %s`, opts.InterfaceName, packagePath))
	gen.SetDefaultParamsPrefix("p")
	gen.SetDefaultResultsPrefix("r")

	v := &visitor{
		gen:             gen,
		methods:         map[string]*types.Signature{},
		sourceInterface: opts.InterfaceName,
	}

	pkg := prog.Package(packagePath)
	if pkg == nil {
		die(fmt.Errorf("unable to load package: %s", packagePath))
	}

	for _, file := range pkg.Files {
		ast.Walk(v, file)
	}

	if len(v.methods) == 0 {
		die(fmt.Errorf("interface %s was not found in %s or it's an empty interface", opts.InterfaceName, packagePath))
	}

	if err := gen.ProcessTemplate("interface", template, v.methods); err != nil {
		die(err)
	}

	if err := gen.WriteToFilename(opts.OutputFile); err != nil {
		die(err)
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch ts := node.(type) {
	case *ast.FuncDecl:
		return nil
	case *ast.TypeSpec:
		exprType, err := v.gen.ExpressionType(ts.Type)
		if err != nil {
			die(fmt.Errorf("failed to get expression for %T %s: %v", ts.Type, ts.Name.Name, err))
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
		input     = flag.String("f", "", "input file or import path of the package that contains interface declaration")
		name      = flag.String("i", "", "name of the interface to mock")
		output    = flag.String("o", "", "destination file name to place the generated mock")
		pkg       = flag.String("p", "", "destination package name")
		sname     = flag.String("t", "", "mock struct name, default is: <interface name>Mock")
		withTests = flag.Bool("withTests", false, "parse *_test.go files in the source package")
	)

	flag.Parse()

	if *pkg == "" || *input == "" || *output == "" || *name == "" || !strings.HasSuffix(*output, ".go") {
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
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
