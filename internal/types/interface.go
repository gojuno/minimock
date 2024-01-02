package types

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
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
	// Find all declared types in a single package
	types := []*ast.TypeSpec{}
	for _, file := range p.Files {
		types = append(types, findAllTypeSpecsInFile(file)...)
	}

	// Filter interfaces from all the declarations
	interfaces := []*ast.TypeSpec{}
	for _, typeSpec := range types {
		if isInterface(typeSpec) {
			interfaces = append(interfaces, typeSpec)
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

func isInterface(typeSpec *ast.TypeSpec) bool {
	// Check if this type declaration is specifically an interface declaration
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
