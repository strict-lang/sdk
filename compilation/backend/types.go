package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
)

const (
	builtinTypeInt    = "int"
	builtinTypeFloat  = "float"
	builtinTypeString = "std::string"
	builtinTypeList   = "std::vector"
)

var builtinTypes = map[string]string{
	"String": builtinTypeString,
	"int":    builtinTypeInt,
	"float":  builtinTypeFloat,
}

func (generation *Generation) GenerateGenericTypeName(name *syntaxtree.GenericTypeName) {
	generation.Emit(name.Name)
	generation.Emit("<")
	generation.EmitNode(name.Generic)
	generation.Emit(">")
}

func (generation *Generation) GenerateListTypeName(name *syntaxtree.ListTypeName) {
	generation.Emit(builtinTypeList)
	generation.Emit("<")
	generation.EmitNode(name.ElementTypeName)
	generation.Emit(">")
}

func (generation *Generation) GenerateConcreteTypeName(name *syntaxtree.ConcreteTypeName) {
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
