package scope

import "gitlab.com/strict-lang/sdk/pkg/compiler/typing"

const builtinScopeId = Id("builtin")

var emptyScope = NewEmptyScope("")

var builtinScope = func() Scope {
	scope := NewOuterScope(builtinScopeId, emptyScope)
	scope.Insert(createNumberType())
	scope.Insert(createFloatType())
	scope.Insert(createBooleanType())
	scope.Insert(createStringType())
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
