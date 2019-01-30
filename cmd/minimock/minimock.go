package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"

	"github.com/marwan-at-work/vgop/modfile"

	"github.com/gojuno/generator"
	"github.com/gojuno/minimock"
)

const version = "1.9.2"

type (
	programOptions struct {
		Interfaces             []interfaceInfo
		Suffix                 string
		OutputFile             string
		StructName             string
		DestinationPackageName string
		ImportWithTests        bool
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
		Methods map[string]*types.Signature
	}

	visitor struct {
		*loader.Program
		interfaces      map[string]interfaceInfo
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
			Error:                    func(err error) {},
		},
	}

	for _, i := range opts.Interfaces {
		if opts.ImportWithTests {
			cfg.ImportWithTests(i.Package)
		} else {
			cfg.Import(i.Package)
		}
	}

	outPackageRealPath := filepath.Dir(opts.OutputFile)
	stat, err := os.Stat(opts.OutputFile)
	if err != nil {
		if !os.IsNotExist(err) {
			die("failed to get file info for %s: %v", opts.OutputFile, err)
		}
	} else if stat.IsDir() {
		outPackageRealPath = opts.OutputFile
	}

	destImportPath, err := packageOf(outPackageRealPath)
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

	packageName := opts.DestinationPackageName
	if packageName == "" {
		packageName = prog.Package(destImportPath).Pkg.Name()
	}

	if len(opts.Interfaces) == 1 && strings.HasSuffix(opts.OutputFile, ".go") { // legacy mode
		interfaceName := opts.Interfaces[0].Name
		sourcePackage := opts.Interfaces[0].Package
		interfaces, err := findInterfaces(prog, interfaceName, sourcePackage)
		if err != nil {
			die("%v", err)
		}

		if len(interfaces) == 0 {
			die("%s was not found in %s", interfaceName, sourcePackage)
		}

		genOpts := generateOptions{
			SourcePackage:      sourcePackage,
			DestinationPackage: destImportPath,
			InterfaceName:      interfaceName,
			StructName:         opts.StructName,
			OutputFileName:     opts.OutputFile,
			PackageName:        packageName,
		}
		if err := generate(prog, genOpts, interfaces[interfaceName].Methods); err != nil {
			die("failed to generate %s: %v", opts.OutputFile, err)
		}
	} else {
		for _, i := range opts.Interfaces {
			interfaces, err := findInterfaces(prog, i.Name, i.Package)
			if err != nil {
				die("%v", err)
			}

			for interfaceName, info := range interfaces {
				genOpts := generateOptions{
					SourcePackage:      i.Package,
					DestinationPackage: destImportPath,
					InterfaceName:      interfaceName,
					StructName:         interfaceName + "Mock",
					OutputFileName:     filepath.Join(outPackageRealPath, minimock.CamelToSnake(interfaceName)+opts.Suffix),
					PackageName:        packageName,
				}

				if err := generate(prog, genOpts, info.Methods); err != nil {
					die("failed to generate %s: %v", genOpts.OutputFileName, err)
				}

				fmt.Printf("Generated file: %s\n", genOpts.OutputFileName)
			}
		}
	}
}

func findInterfaces(prog *loader.Program, sourceInterface, sourcePackage string) (map[string]interfaceInfo, error) {
	v := &visitor{
		Program:         prog,
		sourceInterface: sourceInterface,
		interfaces:      make(map[string]interfaceInfo),
	}

	pkg := prog.Package(sourcePackage)
	if pkg == nil {
		return nil, fmt.Errorf("unable to load package: %s", sourcePackage)
	}

	for _, file := range pkg.Files {
		ast.Walk(v, file)
	}

	return v.interfaces, nil
}

func paramsToStructFields(p generator.ParamSet) string {
	var params []string
	for _, param := range p {
		params = append(params, fmt.Sprintf("%s %s", param.Name, param.Type))
	}

	return strings.Join(params, "\n")
}

func generate(prog *loader.Program, opts generateOptions, methods map[string]*types.Signature) error {
	gen := generator.New(prog)
	gen.ImportWithAlias(opts.DestinationPackage, "")
	gen.SetPackageName(opts.PackageName)
	gen.AddTemplateFunc("toStructFields", paramsToStructFields)
	gen.SetVar("structName", opts.StructName)
	gen.SetVar("interfaceName", opts.InterfaceName)
	gen.SetVar("packagePath", opts.SourcePackage)
	gen.SetHeader(fmt.Sprintf(`DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v%s
The original interface %q can be found in %s`, version, opts.InterfaceName, opts.SourcePackage))
	gen.SetDefaultParamsPrefix("p")
	gen.SetDefaultResultsPrefix("r")

	if len(methods) == 0 {
		return fmt.Errorf("empty interface: %s", opts.InterfaceName)
	}

	if err := gen.ProcessTemplate("interface", template, methods); err != nil {
		return err
	}

	if err := os.Remove(opts.OutputFileName); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove output file %s: %v", opts.OutputFileName, err)
	}

	return gen.WriteToFilename(opts.OutputFileName)
}

// Visit implements ast.Visitor interface
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch ts := node.(type) {
	case *ast.FuncDecl:
		return nil
	case *ast.TypeSpec:
		exprType, err := v.expressionType(ts.Type)
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
		default:
			return nil
		}

		if ts.Name.Name == v.sourceInterface || v.sourceInterface == "*" {
			v.interfaces[ts.Name.Name] = interfaceInfo{
				Name:    ts.Name.Name,
				Methods: getInterfaceMethodsSignatures(i),
			}
		}

		return nil
	}

	return v
}

func (v *visitor) expressionType(e ast.Expr) (types.Type, error) {
	for _, info := range v.Program.AllPackages {
		if typesType := info.TypeOf(e); typesType != nil {
			return typesType, nil
		}
	}

	return nil, fmt.Errorf("expression not found: %+v", e)
}

func getInterfaceMethodsSignatures(t *types.Interface) map[string]*types.Signature {
	methods := make(map[string]*types.Signature)

	for i := 0; i < t.NumMethods(); i++ {
		methods[t.Method(i).Name()] = t.Method(i).Type().(*types.Signature)
	}

	return methods
}

const template = `
	import (
		"github.com/gojuno/minimock"
		testify_assert "github.com/stretchr/testify/assert"
	)

	// {{$structName}} implements {{$packagePath}}.{{$interfaceName}}
	type {{$structName}} struct {
		t minimock.Tester

		{{ range $methodName, $method := . }}
			{{$methodName}}Func func{{ signature $method }}
			{{$methodName}}Counter uint64
			{{$methodName}}PreCounter uint64
			{{$methodName}}Mock m{{$structName}}{{$methodName}}
		{{ end }}
	}

	// New{{$structName}} returns a mock for {{$packagePath}}.{{$interfaceName}}
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
			mock              *{{$structName}}
			mainExpectation   *{{$structName}}{{$methodName}}Expectation
			expectationSeries []*{{$structName}}{{$methodName}}Expectation
		}

		// {{$structName}}{{$methodName}}Expectation specifies expectation struct of the {{$interfaceName}}.{{$methodName}}
		type {{$structName}}{{$methodName}}Expectation struct {
			{{if gt (len (params $method)) 0 }}  input  *{{$structName}}{{$methodName}}Input  {{end}}
			{{if gt (len (results $method)) 0 }} result *{{$structName}}{{$methodName}}Result {{end}}
		}

		{{if gt (len (params $method)) 0 }}
			// {{$structName}}{{$methodName}}Input represents input parameters of the {{$interfaceName}}.{{$methodName}}
			type {{$structName}}{{$methodName}}Input struct {
				{{toStructFields (params $method)}}
			}
		{{end}}

		{{if gt (len (results $method)) 0 }}
			// {{$structName}}{{$methodName}}Result represents results of the {{$interfaceName}}.{{$methodName}}
			type {{$structName}}{{$methodName}}Result struct {
				{{toStructFields (results $method)}}
			}
		{{end}}

		// Expect specifies that invocation of {{$interfaceName}}.{{$methodName}} is expected from 1 to Infinity times
		func (m *m{{$structName}}{{$methodName}}) Expect({{params $method}}) *m{{$structName}}{{$methodName}} {
			m.mock.{{$methodName}}Func = nil
			m.expectationSeries = nil

			if m.mainExpectation == nil {
				m.mainExpectation = &{{$structName}}{{$methodName}}Expectation{}
			}
			{{if gt (len (params $method)) 0 }} m.mainExpectation.input = &{{$structName}}{{$methodName}}Input{ {{ (params $method).Names }} } {{end}}
			return m
		}

		// Return specifies results of invocation of {{$interfaceName}}.{{$methodName}}
		func (m *m{{$structName}}{{$methodName}}) Return({{results $method}}) *{{$structName}} {
			m.mock.{{$methodName}}Func = nil
			m.expectationSeries = nil

			if m.mainExpectation == nil {
				m.mainExpectation = &{{$structName}}{{$methodName}}Expectation{}
			}
			{{if gt (len (results $method)) 0 }} m.mainExpectation.result = &{{$structName}}{{$methodName}}Result{ {{ (results $method).Names }} } {{end}}
			return m.mock
		}

		// ExpectOnce specifies that invocation of {{$interfaceName}}.{{$methodName}} is expected once
		func (m *m{{$structName}}{{$methodName}}) ExpectOnce({{params $method}}) *{{$structName}}{{$methodName}}Expectation {
			m.mock.{{$methodName}}Func = nil
			m.mainExpectation = nil

			expectation := &{{$structName}}{{$methodName}}Expectation{}
			{{if gt (len (params $method)) 0 }} expectation.input = &{{$structName}}{{$methodName}}Input{ {{ (params $method).Names }} } {{end}}
			m.expectationSeries = append(m.expectationSeries, expectation)
			return expectation
		}

		{{if gt (len (results $method)) 0 }}
		// Return sets up return arguments of expectation struct for {{$interfaceName}}.{{$methodName}}
		func (e *{{$structName}}{{$methodName}}Expectation) Return({{results $method}}) {
			e.result = &{{$structName}}{{$methodName}}Result{ {{ (results $method).Names }} }
		}
		{{end}}

		// Set uses given function f as a mock of {{$interfaceName}}.{{$methodName}} method
		func (m *m{{$structName}}{{$methodName}}) Set(f func({{params $method}}) ({{results $method}})) *{{$structName}}{
			m.mainExpectation = nil
			m.expectationSeries = nil

			m.mock.{{$methodName}}Func = f
			return m.mock
		}

		// {{$methodName}} implements {{$packagePath}}.{{$interfaceName}} interface
		func (m *{{$structName}}) {{$methodName}}{{signature $method}} {
			counter := atomic.AddUint64(&m.{{$methodName}}PreCounter, 1)
			defer atomic.AddUint64(&m.{{$methodName}}Counter, 1)

			if len(m.{{$methodName}}Mock.expectationSeries) > 0 {
				if counter > uint64(len(m.{{$methodName}}Mock.expectationSeries)) {
					m.t.Fatalf("Unexpected call to {{$structName}}.{{$methodName}}.{{range (params $method)}} %v{{end}}", {{ (params $method).Names }} )
					return
				}

				{{if gt (len (params $method)) 0 }}
					input := m.{{$methodName}}Mock.expectationSeries[counter-1].input
					testify_assert.Equal(m.t, *input, {{$structName}}{{$methodName}}Input{ {{ (params $method).Names }} }, "{{$interfaceName}}.{{$methodName}} got unexpected parameters")
				{{ end }}

				{{if gt (len (results $method)) 0 }}
 					result := m.{{$methodName}}Mock.expectationSeries[counter-1].result
					if result == nil {
						m.t.Fatal("No results are set for the {{$structName}}.{{$methodName}}")
						return
					}
					{{ range $param := (results $method) }}
					{{ $param.Name }} = result.{{ $param.Name }} {{ end }}
				{{ end }}
				return
			}

			if m.{{$methodName}}Mock.mainExpectation != nil {
				{{if gt (len (params $method)) 0 }}
					input := m.{{$methodName}}Mock.mainExpectation.input
					if input != nil {
						testify_assert.Equal(m.t, *input, {{$structName}}{{$methodName}}Input{ {{ (params $method).Names }} }, "{{$interfaceName}}.{{$methodName}} got unexpected parameters")
					}
				{{ end }}

				{{if gt (len (results $method)) 0 }}
					result := m.{{$methodName}}Mock.mainExpectation.result
					if result == nil {
						m.t.Fatal("No results are set for the {{$structName}}.{{$methodName}}")
					}
					{{ range $param := (results $method) }}
					{{ $param.Name }} = result.{{ $param.Name }} {{ end }}
				{{ end }}
				return
			}

			if m.{{$methodName}}Func == nil {
				m.t.Fatalf("Unexpected call to {{$structName}}.{{$methodName}}.{{range (params $method)}} %v{{end}}", {{ (params $method).Names }} )
				return
			}

			{{if gt (len (results $method)) 0 }} return {{ end }} m.{{$methodName}}Func({{(params $method).Pass}})
		}

		// {{$methodName}}MinimockCounter returns a count of {{$structName}}.{{$methodName}}Func invocations
		func (m *{{$structName}}) {{$methodName}}MinimockCounter() uint64 {
			return atomic.LoadUint64(&m.{{$methodName}}Counter)
		}

		// {{$methodName}}MinimockPreCounter returns the value of {{$structName}}.{{$methodName}} invocations
		func (m *{{$structName}}) {{$methodName}}MinimockPreCounter() uint64 {
			return atomic.LoadUint64(&m.{{$methodName}}PreCounter)
		}

		// {{$methodName}}Finished returns true if mock invocations count is ok
		func (m *{{$structName}}) {{$methodName}}Finished() bool {
			// if expectation series were set then invocations count should be equal to expectations count
			if len(m.{{$methodName}}Mock.expectationSeries) > 0 {
				return atomic.LoadUint64(&m.{{$methodName}}Counter) == uint64(len(m.{{$methodName}}Mock.expectationSeries))
			}

			// if main expectation was set then invocations count should be greater than zero
			if m.{{$methodName}}Mock.mainExpectation != nil {
				return atomic.LoadUint64(&m.{{$methodName}}Counter) > 0
			}

			// if func was set then invocations count should be greater than zero
			if m.{{$methodName}}Func != nil {
				return atomic.LoadUint64(&m.{{$methodName}}Counter) > 0
			}

			return true
		}
	{{ end }}

	// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
	// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) ValidateCallCounters() {
		{{ range $methodName, $method := . }}
			if !m.{{$methodName}}Finished() {
				m.t.Fatal("Expected call to {{$structName}}.{{$methodName}}")
			}
		{{ end }}
	}

	// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
	// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
	func (m *{{$structName}}) CheckMocksCalled() {
		m.Finish()
	}

	// Finish checks that all mocked methods of the interface have been called at least once
	// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
	func (m *{{$structName}}) Finish() {
		m.MinimockFinish()
	}

	// MinimockFinish checks that all mocked methods of the interface have been called at least once
	func (m *{{$structName}}) MinimockFinish() {
		{{ range $methodName, $method := . }}
			if !m.{{$methodName}}Finished()  {
				m.t.Fatal("Expected call to {{$structName}}.{{$methodName}}")
			}
		{{ end }}
	}

	// Wait waits for all mocked methods to be called at least once
	// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
	func (m *{{$structName}}) Wait(timeout time.Duration) {
		m.MinimockWait(timeout)
	}

	// MinimockWait waits for all mocked methods to be called at least once
	// this method is called by minimock.Controller
	func (m *{{$structName}}) MinimockWait(timeout time.Duration) {
		timeoutCh := time.After(timeout)
		for {
			ok := true
			{{ range $methodName, $method := . }}ok = ok && m.{{$methodName}}Finished()
			{{ end }}

			if ok {
				return
			}

			select {
			case <-timeoutCh:
				{{ range $methodName, $method := . }}
					if !m.{{$methodName}}Finished() {
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

	// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
	// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
	func (m *{{$structName}}) AllMocksCalled() bool {
		{{ range $methodName, $method := . }}
			if !m.{{$methodName}}Finished() {
				return false
			}
		{{ end }}

		return true
	}`

func processFlags() *programOptions {
	var (
		input       = flag.String("f", "", "DEPRECATED: input file or import path of the package that contains interface declaration")
		help        = flag.Bool("h", false, "show this help message")
		interfaces  = flag.String("i", "", "comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader, use io.* notation to generate mocks for all interfaces in an io package")
		output      = flag.String("o", "", "destination file name to place the generated mock or path to destination package when multiple interfaces are given")
		packageName = flag.String("p", "", "DEPRECATED: destination package name")
		suffix      = flag.String("s", "_mock_test.go", "output file name suffix which is added to file names when multiple interfaces are given")
		sname       = flag.String("t", "", "DEPRECATED: mock struct name (default <interface name>Mock)")
		v           = flag.Bool("version", false, "show minimock version")
		withTests   = flag.Bool("withTests", false, "parse *_test.go files in the source package")
	)

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *v {
		fmt.Printf("minimock version: %s\n", version)
		os.Exit(0)
	}

	if *input != "" {
		fmt.Printf("minimock: DEPRECATED FLAG: -f\n")
	}

	if *packageName != "" {
		fmt.Printf("minimock: DEPRECATED FLAG: -p\n")
	}

	if *sname != "" {
		fmt.Printf("minimock: DEPRECATED FLAG: -t\n")
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
	for _, interfacePath := range strings.Split(*interfaces, ",") {
		interfacesList = append(interfacesList, getInterfaceInfo(interfacePath))
	}

	return &programOptions{
		Interfaces:             interfacesList,
		OutputFile:             *output,
		StructName:             *sname,
		DestinationPackageName: *packageName,
		ImportWithTests:        *withTests,
		Suffix:                 *suffix,
	}
}

func getInterfaceInfo(interfacePath string) interfaceInfo {
	chunks := strings.Split(interfacePath, ".")
	realPath := strings.Join(chunks[0:len(chunks)-1], ".")
	if len(chunks) < 2 {
		var err error
		if realPath, err = filepath.Abs("."); err != nil {
			die("cannot resolve interface absolute file path", err)
		}
	}
	return interfaceInfo{
		Package: getImportPath(realPath),
		Name:    chunks[len(chunks)-1],
	}
}

func getImportPath(realPath string) string {
	_, err := os.Stat(realPath)
	if err == nil {
		importPath, err := packageOf(realPath)
		if err != nil {
			die("failed to detect import path of the %s: %v", realPath, err)
		}

		return importPath
	}

	return realPath
}

func packageOf(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get fiel %q absolute path: %v", filePath, err)
	}

	f, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open %q: %v", absPath, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file information for %q: %v", absPath, err)
	}

	path := absPath
	if fi.Mode().IsRegular() {
		if !strings.HasSuffix(path, ".go") {
			return "", fmt.Errorf("%q is not a *.go file", absPath)
		}

		path = filepath.Dir(absPath)
	}

	gopath := os.Getenv("GOPATH")
	for _, gp := range filepath.SplitList(gopath) {
		gp, err = filepath.Abs(gp)
		if err != nil {
			continue
		}

		gp += "/src/"

		if strings.HasPrefix(path, gp) {
			return strings.Replace(path, gp, "", 1), nil
		}
	}

	importPath, err := getModulePath(path)
	if err != nil {
		return "", fmt.Errorf("%q is out of GOPATH and %v", path, err)
	}

	return importPath, nil
}

var errNoModFileFound = errors.New("go.mod file is not found in any of the parent directories")

func getModulePath(dir string) (string, error) {
	modFilePath := filepath.Join(dir, "go.mod")

	stat, err := os.Stat(modFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	} else {
		if !stat.IsDir() {
			contents, err := ioutil.ReadFile(modFilePath)
			if err != nil {
				return "", err
			}

			mod, err := modfile.Parse(modFilePath, contents, nil)
			if err != nil {
				return "", err
			}

			return mod.Module.Mod.Path, nil
		}
	}

	var packageName, parentDir string

	filter := func(fi os.FileInfo) bool {
		name := fi.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") && !fi.IsDir()
	}

	fs := token.NewFileSet()
	packages, err := parser.ParseDir(fs, dir, filter, parser.PackageClauseOnly)
	if packages != nil {
		for _, pkg := range packages {
			packageName = pkg.Name
			break
		}
	}

	if packageName == "" {
		parentDir, packageName = filepath.Split(dir)
		if parentDir == dir {
			return "", errNoModFileFound
		}
	}

	parentModPath, err := getModulePath(parentDir)
	if err != nil {
		return "", err
	}

	return parentModPath + "/" + packageName, nil
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "minimock: "+format+"\n", args...)
	os.Exit(1)
}
