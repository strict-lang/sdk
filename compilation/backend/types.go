package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

const (
	builtinTypeInt = "int"
	builtinTypeFloat = "float"
	builtinTypeString  = "std::string"
	builtinTypeList  = "std::vector"
)

var builtinTypes = map[string]string{
	"String":   builtinTypeString,
	"int": builtinTypeInt,
	"float": builtinTypeFloat,
}

func (generation *Generation) GenerateGenericTypeName(name *ast.GenericTypeName) {
	generation.Emit(name.Name)
	generation.Emit("<")
	generation.EmitNode(name.Generic)
	generation.Emit(">")
}

func (generation *Generation) GenerateListTypeName(name *ast.ListTypeName) {
	generation.Emit(builtinTypeList)
	generation.Emit("<")
	generation.EmitNode(name.ElementTypeName)
	generation.Emit(">")
}

func (generation *Generation) GenerateConcreteTypeName(name *ast.ConcreteTypeName) {
	translatedName := lookupTypeName(name.Name)
	generation.Emit(translatedName)
}

func lookupTypeName(name string) string {
	builtin, ok := builtinTypes[name]
	if !ok {
		return name
	}
	return builtin
}
