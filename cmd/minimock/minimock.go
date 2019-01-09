package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gojuno/minimock"
	"github.com/hexdigest/gowrap/generator"
	"github.com/hexdigest/gowrap/pkg"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

var version = "dev" //do not modify! version var is modified during the build via ldflags option

type (
	options struct {
		interfaces []interfaceInfo
		noGenerate bool
	}

	interfaceInfo struct {
		Type       string
		ImportPath string
		WriteTo    string
	}
)

func main() {
	opts, err := processFlags(os.Args[1:], os.Stderr)
	if err != nil {
		die(2, "%v", err)
	}

	if opts == nil {
		os.Exit(2)
	}

	if err = processOptions(opts); err != nil {
		die(1, "%v", err)
	}
}

func processOptions(opts *options) (err error) {
	var (
		sourcePackage *packages.Package
		astPackage    *ast.Package
		fs            *token.FileSet
	)

	for i, in := range opts.interfaces {
		if i == 0 || in.ImportPath != opts.interfaces[i-1].ImportPath {
			if sourcePackage, err = pkg.Load(in.ImportPath); err != nil {
				return err
			}

			fs = token.NewFileSet()
			if astPackage, err = pkg.AST(fs, sourcePackage); err != nil {
				return errors.Wrap(err, "failed to load package sources")
			}
		}

		interfaces, err := findInterfaces(astPackage, in.Type)
		if err != nil {
			return err
		}

		gopts := generator.Options{
			SourcePackage:  sourcePackage.PkgPath,
			HeaderTemplate: minimock.HeaderTemplate,
			BodyTemplate:   minimock.BodyTemplate,
			HeaderVars: map[string]interface{}{
				"GenerateInstruction": !opts.noGenerate,
				"Version":             version,
			},
			Funcs: template.FuncMap{
				"title": strings.Title,
			},
		}

		if err := processPackage(gopts, interfaces, in.WriteTo); err != nil {
			return err
		}
	}

	return nil
}

func processPackage(opts generator.Options, interfaces []string, writeTo string) (err error) {
	for _, name := range interfaces {
		opts.InterfaceName = name
		opts.OutputFile, err = destinationFile(name, writeTo)
		if err != nil {
			return errors.Wrapf(err, "failed to generate mock for %s", name)
		}

		if err := generate(opts); err != nil {
			return err
		}

		fmt.Printf("minimock: %s\n", opts.OutputFile)
	}

	return nil
}

func isGoFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}

		dir := filepath.Dir(path)
		if stat, err = os.Stat(dir); err != nil {
			return false, err
		}

		if !stat.IsDir() {
			return false, errors.Errorf("not a directory: %s", dir)
		}

		return strings.HasSuffix(path, ".go"), nil
	}

	return strings.HasSuffix(path, ".go") && !stat.IsDir(), nil
}

func destinationFile(interfaceName, writeTo string) (string, error) {
	ok, err := isGoFile(writeTo)
	if err != nil {
		return "", err
	}

	var path string

	if ok {
		path = writeTo
	} else {
		path = filepath.Join(writeTo, minimock.CamelToSnake(interfaceName)+"_mock_test.go")
	}

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, ".") {
		path = "./" + path
	}

	return path, nil
}

func generate(o generator.Options) (err error) {
	g, err := generator.NewGenerator(o)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})

	if err = g.Generate(buf); err != nil {
		return errors.Wrap(err, "failed to generate mock")
	}

	return ioutil.WriteFile(o.OutputFile, buf.Bytes(), 0644)
}

func findInterfaces(p *ast.Package, pattern string) ([]string, error) {
	var names []string
	for _, f := range p.Files {
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
				for _, spec := range gd.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := ts.Type.(*ast.InterfaceType); ok && match(ts.Name.Name, pattern) {
							names = append(names, ts.Name.Name)
						}
					}
				}
			}
		}
	}

	if len(names) == 0 {
		return nil, errors.Errorf("failed to find any interfaces matching %s in %s", pattern, p.Name)
	}

	return names, nil
}

func match(s, pattern string) bool {
	return pattern == "*" || s == pattern
}

func processFlags(args []string, stderr io.Writer) (*options, error) {
	var opts options

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&opts.noGenerate, "g", false, "don't put go:generate instruction into the generated code")
	input := fs.String("i", "*", "comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader\nuse io.* notation to generate mocks for all interfaces in the io package")
	output := fs.String("o", "", "comma-separated destination file name(s) or package(s) to put the generated mocks in,\nby default the generated mock is placed into the source package directory")

	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: \033[1mminimock\033[0m [\033[1m-i\033[0m source.interface] [\033[1m-o\033[0m output/dir/or/file.go] [\033[1m-g\033[0m]")
		fs.PrintDefaults()

		fmt.Fprintf(stderr, "\nExamples:\n\n")
		fmt.Fprintf(stderr, "  Generate mocks for all interfaces that can be found in the current directory:\n")
		fmt.Fprintf(stderr, "    \033[1mminimock\033[0m\n\n")

		fmt.Fprintf(stderr, "  Generate mock for the io.Writer interface and put it into the \"./buffer\" subpackage/dir:\n")
		fmt.Fprintf(stderr, "    \033[1mminimock\033[0m \033[1m-i\033[0m io.Writer \033[1m-o\033[0m ./buffer\n\n")

		fmt.Fprintf(stderr, "  Generate mocks for the fmt.Stringer and all interfaces from the \"io\" package and put them into the \"./buffer\" subpackage/dir:\n")
		fmt.Fprintf(stderr, "    \033[1mminimock\033[0m \033[1m-i\033[0m fmt.Stringer,io.* \033[1m-o\033[0m ./buffer\n\n")

		fmt.Fprintln(stderr, "\nFor more information please visit https://github.com/gojuno/minimock")
	}

	if err := fs.Parse(args); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(stderr, "%v\n\n", err)
			fs.Usage()
		}
		return nil, nil
	}

	interfaces := strings.Split(*input, ",")

	var writeTo = make([]string, len(interfaces))
	if *output != "" {
		//if only one output package/file was given
		if to := strings.Split(*output, ","); len(to) == 1 {
			for i := range writeTo {
				writeTo[i] = to[0]
			}
		} else {
			writeTo = to
		}
	}

	if len(writeTo) != len(interfaces) {
		return nil, errors.Errorf("count of the source interfaces doesn't match the output files count")
	}

	if err := checkDuplicateOutputFiles(writeTo); err != nil {
		return nil, err
	}

	for i, in := range interfaces {
		info, err := makeInterfaceInfo(in, writeTo[i])
		if err != nil {
			return nil, err
		}

		opts.interfaces = append(opts.interfaces, *info)
	}

	return &opts, nil
}

// checkDuplicateOutputFiles finds first non-unique Go file
func checkDuplicateOutputFiles(fileNames []string) error {
	for i := range fileNames {
		ok, err := isGoFile(fileNames[i])
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		ipath, err := filepath.Abs(fileNames[i])
		if err != nil {
			return err
		}

		for j := range fileNames {
			jpath, err := filepath.Abs(fileNames[j])
			if err != nil {
				return err
			}

			if j != i && ipath == jpath {
				return errors.Errorf("duplicate output file name: %s", ipath)
			}
		}
	}

	return nil
}

func makeInterfaceInfo(typ, writeTo string) (*interfaceInfo, error) {
	info := interfaceInfo{WriteTo: writeTo}

	dot := strings.LastIndex(typ, ".")
	slash := strings.LastIndex(typ, "/")

	if slash > dot {
		return nil, errors.Errorf("invalid interface type: %s", typ)
	}

	if dot >= 0 {
		info.Type = typ[dot+1:]
		info.ImportPath = typ[:dot]
	} else {
		info.Type = typ
		info.ImportPath = "./"
	}

	return &info, nil
}

func die(exitCode int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "minimock: "+format+"\n", args...)
	os.Exit(exitCode)
}
