package code



var builtinTypes = struct {
	boolType *Type
	intType *Type
	floatType *Type
	invalidType *Type
	stringType *Type
} {
	boolType: createPrimitiveType("bool"),
	intType: createPrimitiveType("int"),
	floatType: createPrimitiveType("float"),
	invalidType: createPrimitiveType("invalid"),
	stringType: createPrimitiveType("string"),
}

func createPrimitiveType(name string) *Type {
	return &Type{
		Name:                      name,
		Methods:                   map[string]*TypedMethod{},
		Descriptor:                TypeDescriptor(name),
		ParameterCount:            0,
		ImplicitConversionTargets: []*Type{},
		SuperType:                 nil, // TODO(): Common super type
		SuperInterfaces:           []*Type{},
	}
}
