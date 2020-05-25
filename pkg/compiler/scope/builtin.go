package scope

import "github.com/strict-lang/sdk/pkg/compiler/typing"

const builtinScopeId = Id("builtin")

var emptyScope = NewEmptyScope("")
var builtinScope = NewOuterScopeWithRootId(builtinScopeId, emptyScope)
var booleanType = createBooleanType()

var Builtins = struct {
	Any     *Class
	Number  *Class
	Float   *Class
	Boolean *Class
	String  *Class
	Void    *Class
	Computation *Class
	True    *Field
	False   *Field
}{
	Void:    createVoidType(),
	Number:  createNumberType(),
	Float:   createFloatType(),
	Any:     createAnyType(),
	Boolean: booleanType,
	String:  createStringType(),
	True:    createBuiltinField("True", booleanType),
	False:   createBuiltinField("False", booleanType),
	Computation: createComputationType(),
}

func init() {
	builtinScope.Insert(Builtins.Number)
	builtinScope.Insert(Builtins.Float)
	builtinScope.Insert(Builtins.Boolean)
	builtinScope.Insert(Builtins.String)
	builtinScope.Insert(Builtins.True)
	builtinScope.Insert(Builtins.False)
	builtinScope.Insert(Builtins.Void)
	builtinScope.Insert(Builtins.Any)
}

func createAnyType() *Class {
	any := createPrimitiveClass("Any")
	any.ActualClass = typing.NewEmptyClass("Any")
	any.Scope = createAnyContents()
	return any
}

func createAnyContents() MutableScope {
	contents := NewOuterScope("Builtin.Any", emptyScope)
	contents.Insert(&Method{
		DeclarationName:   "CalculateHashCode",
		declarationOffset: 0,
		ReturnType:        booleanType,
	})
	return contents
}

func createComputationType() *Class {
	actualClass := typing.NewEmptyClass("Computation")
	return &Class{
		DeclarationName:   "Computation",
		QualifiedName:     "Strict.Base.Computation",
		ActualClass:       actualClass,
	}
}

func createVoidType() *Class {
	number := createPrimitiveClass("Void")
	number.ActualClass = typing.NewEmptyClass("Void")
	return number
}

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

func createBuiltinField(name string, class *Class) *Field {
	return &Field{
		DeclarationName:   name,
		declarationOffset: 0,
		Class:             class,
		Kind:              ConstantField,
	}
}

func createPrimitiveClass(name string) *Class {
	return &Class{
		DeclarationName: name,
		Scope:           NewOuterScope(Id(name), builtinScope),
	}
}

func NewBuiltinScope() Scope {
	return builtinScope
}
