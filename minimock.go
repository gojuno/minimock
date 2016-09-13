package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type (
	options struct {
		InputFile     string
		OutputFile    string
		InterfaceName string
		StructName    string
		Package       string
	}

	generator struct {
		options
		content      *bytes.Buffer
		fileset      *token.FileSet
		err          error
		importsPaths map[string]string
		importsUsage map[string]bool
		pkg          string
	}
)

func main() {
	opts := processFlags()
	if opts == nil {
		os.Exit(1)
	}

	g := newGenerator(opts)

	generatedContent, err := g.generate()
	if err != nil {
		die(err)
	}

	f, err := os.Create(opts.OutputFile)
	if err != nil {
		die(err)
	}

	if _, err := f.Write(generatedContent); err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func processFlags() *options {
	var (
		input  = flag.String("f", "", "input file containing interface declaration")
		name   = flag.String("i", "", "interface name")
		output = flag.String("o", "", "destination file for interface implementation")
		pkg    = flag.String("p", "", "destination package name")
		sname  = flag.String("t", "", "target struct name, default: <interface name>Mock")
	)

	flag.Parse()

	if *pkg == "" || *input == "" || *output == "" || *name == "" || !strings.HasSuffix(*input, ".go") || !strings.HasSuffix(*output, ".go") {
		flag.Usage()
		return nil
	}

	if *sname == "" {
		*sname = *name + "Mock"
	}

	return &options{
		InputFile:     *input,
		OutputFile:    *output,
		InterfaceName: *name,
		Package:       *pkg,
		StructName:    *sname,
	}
}

func newGenerator(o *options) *generator {
	return &generator{
		options:      *o,
		content:      new(bytes.Buffer),
		importsPaths: map[string]string{},
		importsUsage: map[string]bool{},
	}
}

func (g *generator) generate() ([]byte, error) {
	g.fileset = token.NewFileSet()

	parsedFile, err := parser.ParseFile(g.fileset, g.InputFile, nil, 0)
	if err != nil {
		return nil, err
	}

	g.pkg = parsedFile.Name.Name

	inputFilePath, err := filepath.Abs(g.InputFile)
	if err != nil {
		return nil, err
	}

	outputFilePath, err := filepath.Abs(g.OutputFile)
	if err != nil {
		return nil, err
	}

	if g.pkg == "main" && inputFilePath != outputFilePath {
		return nil, fmt.Errorf("file %s is a part of the main package, output file should be in the same dir", g.InputFile)
	}

	g.importsPaths = processImports(parsedFile)

	//we have selector in input file that is equal to input file's package name
	if _, ok := g.importsPaths[g.pkg]; ok {
		g.pkg = "_root_pkg"
	}

	g.importsPaths[g.pkg] = getRelativePkgPath(inputFilePath)
	ast.Inspect(parsedFile, g.interfaceInspector)

	if g.err != nil {
		return nil, g.err
	}

	return g.getGeneratedContent()
}

func (g *generator) interfaceInspector(node ast.Node) bool {
	typeSpec, ok := node.(*ast.TypeSpec)
	if ok && typeSpec.Name.Name == g.InterfaceName {
		if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
			g.err = fmt.Errorf("%s: type %s is not an interface", g.InputFile, g.InterfaceName)
		} else {
			g.processInterfaceDeclaration(interfaceType)
		}
	}

	return !ok
}

const (
	templateForHeader = `
		package {{package}}

		import "testing"
		import "sync"
		{{imports}}

		type {{struct}} struct {
			t *testing.T
			m *sync.Mutex

			{{methods}}

			{{counters}}
		}

		func New{{struct}}(t *testing.T) *{{struct}} {
			return &{{struct}}{t: t, m: &sync.Mutex{}}
		}
		`

	templateForMethod = `
		func (m *{{struct}}) {{method}}{{inParams}} {{outParams}} {
			m.m.Lock()
			m.{{method}}Counter += 1
			m.m.Unlock()

			if m.{{method}}Func == nil {
				m.t.Fatalf("Unexpected call to {{struct}}.{{method}}")
			}

			{{return}}m.{{method}}Func({{paramsList}})
		}
		`
	templateValidateHeader = `
		func (m *{{struct}}) ValidateCallCounters() {
	`

	templateValidateBody = `
		if m.{{method}}Func != nil && m.{{method}}Counter == 0 {
			m.t.Error("Expected call to {{struct}}.{{method}}")
		}

	`
	templateValidateFooter = `
		}
	`
)

func (g *generator) processTemplate(template string, vars map[string]string) {
	for k, v := range vars {
		template = strings.Replace(template, "{{"+k+"}}", v, -1)
	}

	fmt.Fprintf(g.content, template)
}

func (g *generator) processInterfaceDeclaration(it *ast.InterfaceType) error {
	methods := []string{}
	counters := []string{}
	names := []string{}

	for i, method := range it.Methods.List {
		methodName := it.Methods.List[i].Names[0].Name
		names = append(names, methodName)

		signature := fmt.Sprintf("  %sFunc %s", methodName, g.getMethodSignature(method.Type))

		methods = append(methods, signature)
		counters = append(counters, methodName+"Counter int")

		f := method.Type.(*ast.FuncType)

		g.processUsedImports(f.Params)
		g.processUsedImports(f.Results)
	}

	vars := map[string]string{
		"struct":   g.StructName,
		"package":  g.Package,
		"methods":  strings.Join(methods, "\n"),
		"counters": strings.Join(counters, "\n"),
		"imports":  g.getImports(),
	}

	g.processTemplate(templateForHeader, vars)

	for _, method := range it.Methods.List {
		f := method.Type.(*ast.FuncType)

		vars["method"] = method.Names[0].Name
		vars["inParams"] = g.getInParams(f)
		vars["outParams"] = g.getOutParams(f)

		vars["paramsList"] = strings.Join(g.getParamsList(f.Params.List), ", ")
		if f.Results == nil {
			vars["return"] = ""
		} else {
			vars["return"] = "return "
		}

		g.processTemplate(templateForMethod, vars)
	}

	g.processTemplate(templateValidateHeader, vars)
	for _, name := range names {
		vars["method"] = name
		g.processTemplate(templateValidateBody, vars)
	}
	g.processTemplate(templateValidateFooter, vars)

	return nil
}

func processImports(f *ast.File) map[string]string {
	result := map[string]string{}

	for _, imp := range f.Imports {
		selector := ""
		path := imp.Path.Value[1 : len(imp.Path.Value)-1]

		if imp.Name == nil {
			selector = filepath.Base(path)
		} else {
			selector = strings.SplitN(imp.Name.Name, ".", 1)[0]
		}

		result[selector] = path
	}

	return result
}

func (g *generator) getImports() string {
	imports := []string{}

	for selector, _ := range g.importsUsage {
		path := g.importsPaths[selector]
		importSpec := ""
		if path == selector || strings.HasSuffix(path, "/"+selector) {
			importSpec = fmt.Sprintf("import %q", path)
		} else {
			importSpec = fmt.Sprintf("import %s %q", selector, path)
		}

		imports = append(imports, importSpec)
	}

	return strings.Join(imports, "\n")
}

func (g *generator) getParamsList(fields []*ast.Field) []string {
	result := []string{}

	for i, f := range fields {
		name := ""
		if len(f.Names) == 0 {
			name = fmt.Sprintf("p%d_0%s", i, variadicSuffix(f))
			result = append(result, name)
		} else {
			for _, n := range f.Names {
				name = n.Name + variadicSuffix(f)
				result = append(result, name)
			}
		}
	}

	return result
}

func (g *generator) processUsedImports(fields *ast.FieldList) {
	if fields == nil {
		return
	}

	for _, f := range fields.List {
		if selector := g.getSelector(f.Type); selector != "" {
			g.importsUsage[selector] = true
		}
	}
}

func (g *generator) getInParams(f *ast.FuncType) string {
	params := []string{}

	for i, f := range f.Params.List {
		names := []string{}
		if len(f.Names) == 0 {
			names = append(names, fmt.Sprintf("p%d_0", i))
		} else {
			for _, n := range f.Names {
				names = append(names, n.Name)
			}
		}

		params = append(params, strings.Join(names, ", ")+" "+g.typeToString(f.Type))
	}

	return "(" + strings.Join(params, ", ") + ")"
}

func (g *generator) getOutParams(f *ast.FuncType) string {
	if f.Results == nil {
		return ""
	}

	params := []string{}

	for _, f := range f.Results.List {
		types := []string{}
		if len(f.Names) == 0 {
			types = append(types, variadicSuffix(f)+g.typeToString(f.Type))
		} else {
			for range f.Names {
				types = append(types, variadicSuffix(f)+g.typeToString(f.Type))
			}
		}

		params = append(params, strings.Join(types, ", "))
	}

	if len(params) < 2 {
		return strings.Join(params, ", ")
	}

	return "(" + strings.Join(params, ", ") + ")"
}

func (g *generator) typeToString(t ast.Expr) string {
	switch e := t.(type) {
	case *ast.StarExpr:
		return "*" + g.typeToString(e.X)
	case *ast.ArrayType:
		return "[]" + g.typeToString(e.Elt)
	case *ast.MapType:
		return "map[" + g.typeToString(e.Key) + "]" + g.typeToString(e.Value)
	case *ast.SelectorExpr:
		return g.getSelector(e) + "." + e.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		switch e.Dir {
		case ast.SEND:
			return "chan <- " + g.typeToString(e.Value)
		case ast.RECV:
			return "<- chan " + g.typeToString(e.Value)
		default:
			return "chan " + g.typeToString(e.Value)
		}
	case *ast.Ellipsis:
		return "..." + g.typeToString(e.Elt)
	case *ast.Ident:
		name := e.Name
		if strings.ToUpper(name)[0] == name[0] {
			return g.pkg + "." + name
		} else {
			return name
		}
	}

	panic(fmt.Sprintf("Unexpected expression of type %T: %+v", t, t))
}

func (g *generator) getSelector(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name
	case *ast.ArrayType:
		return g.getSelector(t.Elt)
	case *ast.StarExpr:
		return g.getSelector(t.X)
	case *ast.Ident:
		name := t.Name
		if strings.ToUpper(name)[0] == name[0] {
			return g.pkg
		}
	}

	return ""
}

func (g *generator) getGeneratedContent() ([]byte, error) {
	formatted, err := format.Source(g.content.Bytes())
	if err != nil {
		fmt.Printf(string(g.content.Bytes()))
		return nil, err
	}

	return formatted, nil
}

func (g *generator) getMethodSignature(t ast.Expr) string {
	return "func" + g.getInParams(t.(*ast.FuncType)) + " " + g.getOutParams(t.(*ast.FuncType))
}

func getRelativePkgPath(inputFilePath string) string {
	pkgDir := filepath.Dir(inputFilePath)

	paths := strings.Split(os.Getenv("GOPATH"), ":")
	for _, p := range paths {
		prefix := p + "/src/"
		if strings.HasPrefix(pkgDir, prefix) {
			return strings.Replace(pkgDir, prefix, "", 1)
		}
	}

	return pkgDir
}

func variadicSuffix(f *ast.Field) string {
	if _, ok := f.Type.(*ast.Ellipsis); ok {
		return "..."
	}

	return ""
}
