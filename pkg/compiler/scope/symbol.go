package scope

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/typing"
)

type Symbol interface {
	Name() string
	String() string
	DeclarationOffset() input.Offset
}

type Method struct {
	DeclarationName   string
	declarationOffset input.Offset
	ReturnType *Class
	// Parameters are lazily added
	Parameters []*Field
}

func (method *Method) Name() string {
	return method.DeclarationName
}

func (method *Method) String() string {
	return fmt.Sprintf("Method{Name: %s, ReturnType: %s}", method.DeclarationName, method.ReturnType)
}

func (method *Method) DeclarationOffset() input.Offset {
	return method.declarationOffset
}

type Class struct {
	DeclarationName   string
	ActualClass       typing.Type
	declarationOffset input.Offset
}

func (class *Class) Name() string {
	return class.DeclarationName
}

func (class *Class) String() string {
	return fmt.Sprintf("Class{Name: %s}", class.DeclarationName)
}

func (class *Class) DeclarationOffset() input.Offset {
	return class.declarationOffset
}

type Field struct {
	DeclarationName string
	declarationOffset input.Offset
	Class *Class
	Kind FieldKind
}

type FieldKind int

const (
	ParameterField FieldKind = iota
	VariableField
	MemberField
	ConstantField
)

func (field *Field) Name() string {
	return field.DeclarationName
}

func (field *Field) String() string {
	return fmt.Sprintf("Field{Name: %s, Type: %s}", field.DeclarationName, field.Class)
}

func (field *Field) DeclarationOffset() input.Offset {
	return field.declarationOffset
}

type Namespace struct {
	PackageName string
	Scope Scope
}

func (namespace *Namespace) Name() string {
	return namespace.PackageName
}

func (namespace *Namespace) DeclarationOffset() input.Offset {
	return 0
}

func (namespace *Namespace) String() string {
	return fmt.Sprintf("Namespace{Name: %s}", namespace.PackageName)
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

func AsNamespaceSymbol(symbol Symbol) (*Namespace, bool) {
	namespace, ok := symbol.(*Namespace)
	return namespace, ok
}
