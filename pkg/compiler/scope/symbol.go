package scope

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/typing"
)

type Symbol interface {
	Name() string
	DeclarationOffset() input.Offset
}

type Method struct {
	name string
	declarationOffset input.Offset

	ReturnType typing.Class
	Parameters []*Field
}

func (method *Method) Name() string {
	return method.name
}

func (method *Method) DeclarationOffset() input.Offset {
	return method.declarationOffset
}

type Class struct {
	name string
	class typing.Class
	declarationOffset input.Offset
}

func (class *Class) Name() string {
	return class.name
}

func (class *Class) DeclarationOffset() input.Offset {
	return class.declarationOffset
}

type Field struct {
	name string
	declarationOffset input.Offset
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) DeclarationOffset() input.Offset {
	return field.declarationOffset
}

func AsMethodSymbol(symbol Symbol) (*Method, bool) {
	method, ok := symbol.(*Method)
	return method, ok
}

func AsClassSymbol(symbol Symbol) (*Class, bool) {
	class, ok := symbol.(*Class)
	return class, ok
}

func AsFieldSymbol(symbol Symbol) (*Field, bool) {
	field, ok := symbol.(*Field)
	return field, ok
}
