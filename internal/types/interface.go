package types

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/hexdigest/gowrap/pkg"
)

// InterfaceSpecification represents abstraction over interface type. It contains all the metadata
// required to render a mock for given interface. One could deduce whether interface is generic
// by looking for type params
type InterfaceSpecification struct {
	InterfaceName   string
	InterfaceParams []InterfaceSpecificationParam
}

// InterfaceSpecificationParam represents a group of type param variables and their type
// I.e. [T,K any] would result in names "T","K" and type "any"
type InterfaceSpecificationParam struct {
	ParamNames []string
	ParamType  string
}

func FindAllInterfaces(p *ast.Package, pattern string) []InterfaceSpecification {
	// Filter interfaces from all the declarations
	interfaces := []*ast.TypeSpec{}
	for _, file := range p.Files {
		for _, typeSpec := range findAllTypeSpecsInFile(file) {
			if isInterface(typeSpec, file.Imports) {
				interfaces = append(interfaces, typeSpec)
			}
		}
	}

	// Filter interfaces with the given pattern
	filteredInterfaces := []*ast.TypeSpec{}
	for _, iface := range interfaces {
		if match(iface.Name.Name, pattern) {
			filteredInterfaces = append(filteredInterfaces, iface)
		}
	}

	// Transform AST nodes into specifications
	interfaceSpecifications := make([]InterfaceSpecification, 0, len(filteredInterfaces))
	for _, iface := range filteredInterfaces {
		interfaceSpecifications = append(interfaceSpecifications, InterfaceSpecification{
			InterfaceName:   iface.Name.Name,
			InterfaceParams: getTypeParams(iface),
		})
	}

	return interfaceSpecifications
}

func isInterface(typeSpec *ast.TypeSpec, fileImports []*ast.ImportSpec) bool {
	// we are generating mocks for interfaces,
	// interface aliases to types from the same package
	// and aliases to types from another package
	return isInterfaceType(typeSpec) ||
		isInterfaceAlias(typeSpec, fileImports) ||
		isExportedInterfaceAlias(typeSpec, fileImports)
}

// isInterfaceAlias checks if type is an alias to other
// interface type that is in the same package
func isInterfaceAlias(typeSpec *ast.TypeSpec, fileImports []*ast.ImportSpec) bool {
	ident, ok := typeSpec.Type.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Obj == nil {
		return false
	}

	if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
		return isInterface(ts, fileImports)
	}

	return false
}

// isExportedInterfaceAlias checks if type is an alias to other
// interface type that is in other exported package
func isExportedInterfaceAlias(typeSpec *ast.TypeSpec, fileImports []*ast.ImportSpec) bool {
	selector, ok := typeSpec.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	name := selector.Sel.Name
	srcPkgPath := findSourcePackage(ident, fileImports)
	srcPackageAst, err := getPackageAst(srcPkgPath)
	if err != nil {
		return false
	}

	for _, f := range srcPackageAst.Files {
		if f == nil {
			continue
		}
		types := findAllTypeSpecsInFile(f)
		typeSpec, found := findTypeByName(types, name)
		if found {
			// we have to check recursively because checked typed might be
			// another alias to other interface
			return isInterface(typeSpec, f.Imports)
		}
	}

	return false
}

func getPackageAst(packagePath string) (*ast.Package, error) {
	srcPkg, err := pkg.Load(packagePath)
	if err != nil {
		return nil, err
	}

	fs := token.NewFileSet()
	srcAst, err := pkg.AST(fs, srcPkg)
	if err != nil {
		return nil, err
	}

	return srcAst, nil
}

func findTypeByName(types []*ast.TypeSpec, name string) (*ast.TypeSpec, bool) {
	for _, ts := range types {
		if ts.Name.Name == name {
			return ts, true
		}
	}

	return nil, false
}

func findSourcePackage(ident *ast.Ident, imports []*ast.ImportSpec) string {
	for _, imp := range imports {
		cleanPath := strings.Trim(imp.Path.Value, "\"")
		if imp.Name != nil {
			if ident.Name == imp.Name.Name {
				return cleanPath
			}

			continue
		}

		// try last segment, like in "github.com/my/package/identName"
		lastSlash := strings.LastIndex(cleanPath, "/")
		if ident.Name == cleanPath[lastSlash+1:] {
			return cleanPath
		}

		// try prev segment, like in "github.com/my/package/identName/v5"
		if cleanPath[lastSlash+1] == 'v' {
			prevSlash := strings.LastIndex(cleanPath[:lastSlash], "/")
			if ident.Name == cleanPath[prevSlash+1:lastSlash] {
				return cleanPath
			}
		}
	}

	// todo: should not reach here?
	return ""
}

func isInterfaceType(typeSpec *ast.TypeSpec) bool {
	_, ok := typeSpec.Type.(*ast.InterfaceType)
	return ok
}

// findAllInterfaceNodesInFile ranges over file's AST nodes and extracts all interfaces inside
// returned *ast.TypeSpecs can be safely interpreted as interface declaration nodes
func findAllTypeSpecsInFile(f *ast.File) []*ast.TypeSpec {
	typeSpecs := []*ast.TypeSpec{}

	// Range over all declarations in a single file
	for _, declaration := range f.Decls {
		// Check if declaration is an import, constant, type or variable declaration.
		// If it is, check specifically if it's a TYPE as all interfaces are types
		if genericDeclaration, ok := declaration.(*ast.GenDecl); ok && genericDeclaration.Tok == token.TYPE {
			// Range over all specifications and find ones that are Type declarations
			// This is mostly a precaution
			for _, spec := range genericDeclaration.Specs {
				// Check directly for a type spec declaration
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					typeSpecs = append(typeSpecs, typeSpec)
				}
			}
		}
	}

	return typeSpecs
}

// match returns true if pattern is a wildcard or directly matches the given name
func match(name, pattern string) bool {
	return pattern == "*" || name == pattern
}

func getTypeParams(typeSpec *ast.TypeSpec) []InterfaceSpecificationParam {
	params := []InterfaceSpecificationParam{}

	// Check whether node has any type params at all
	if typeSpec == nil || typeSpec.TypeParams == nil {
		return nil
	}

	// If node has any type params - store them in slice and return as a spec
	for _, param := range typeSpec.TypeParams.List {
		names := []string{}
		for _, name := range param.Names {
			names = append(names, name.Name)
		}

		var out bytes.Buffer
		printer.Fprint(&out, token.NewFileSet(), param.Type)

		paramType := out.String()

		params = append(params, InterfaceSpecificationParam{
			ParamNames: names,
			ParamType:  paramType,
		})
	}

	return params
}
