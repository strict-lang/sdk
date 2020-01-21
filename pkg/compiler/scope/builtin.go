package scope

import "gitlab.com/strict-lang/sdk/pkg/compiler/typing"

const builtinScopeId = Id("builtin")

var emptyScope = NewEmptyScope("")

var Builtins = struct {
	Number *Class
	Float *Class
	Boolean *Class
	String *Class
}{
	Number: createNumberType(),
	Float: createFloatType(),
	Boolean: createBooleanType(),
	String: createStringType(),
}

var builtinScope = func() Scope {
	scope := NewOuterScope(builtinScopeId, emptyScope)
	scope.Insert(Builtins.Number)
	scope.Insert(Builtins.Float)
	scope.Insert(Builtins.Boolean)
	scope.Insert(Builtins.String)
	return scope
}()

func createNumberType() *Class {
	number := createPrimitiveClass("Number")
	number.ActualClass = typing.NewEmptyClass("Number")
	return number
}

func createFloatType() *Class {
	class := createPrimitiveClass("Float")
	class.ActualClass = typing.NewEmptyClass("Float")
	return class
}

func createBooleanType() *Class {
	class := createPrimitiveClass("Boolean")
	class.ActualClass = typing.NewEmptyClass("Boolean")
	return class
}

func createStringType() *Class {
	class := createPrimitiveClass("String")
	class.ActualClass = typing.NewEmptyClass("String")
	return class
}

func createPrimitiveClass(name string) *Class {
	return &Class{
		DeclarationName:   name,
	}
}

func NewBuiltinScope() Scope {
	return builtinScope
}
