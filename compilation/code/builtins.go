package code

import "gitlab.com/strict-lang/sdk/compilation/scope"

type BuiltinRegistration struct {
	scope *scope.Scope
}

// AnyType and VoidType are not created using the TypeBuilder since they are
// special essential types that do not have the default parameters for things
// like implicit conversion checks.
var (
	AnyType = Type{
		Name:       "any",
		Descriptor: typeDescriptorAny,
	}

	VoidType = Type{
		Name:       "void",
		Descriptor: typeDescriptorVoid,
	}

	BoolType = TypeBuilder{
		Name:       "bool",
		Descriptor: typeDescriptorBool,
	}.Create()

	NumberType = TypeBuilder{
		Name:       "number",
		Descriptor: typeDescriptorNumber,
	}.Create()

	TextType = TypeBuilder{
		Name:       "text",
		Descriptor: typeDescriptorText,
		Methods: []TypedMethod{
			{
				Name:       "size",
				ReturnType: &NumberType,
			},
			{
				Name:       "isEmpty",
				ReturnType: &BoolType,
			},
		},
	}.Create()
)

func (builtins *BuiltinRegistration) register() {
	builtins.registerType(&AnyType)
	builtins.registerType(&VoidType)
	builtins.registerType(&TextType)
	builtins.registerType(&NumberType)
	builtins.registerType(&BoolType)
}

func (builtins *BuiltinRegistration) registerType(builtinType *Type) {
}

func (builtins *BuiltinRegistration) registerMethod(method *TypedMethod) {}
