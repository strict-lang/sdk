package code

var builtinTypes = struct {
	boolType    *Class
	intType     *Class
	floatType   *Class
	invalidType *Class
	stringType  *Class
}{
	boolType:    createPrimitiveType("bool"),
	intType:     createPrimitiveType("int"),
	floatType:   createPrimitiveType("float"),
	invalidType: createPrimitiveType("invalid"),
	stringType:  createPrimitiveType("string"),
}

func createPrimitiveType(name string) *Class {
	return &Class{
		Name:                      name,
		Methods:                   map[string]*TypedMethod{},
		Descriptor:                TypeDescriptor(name),
		ParameterCount:            0,
		ImplicitConversionTargets: []TypeDescriptor{},
		SuperType:                 nil, // TODO(): Common super type
		SuperInterfaces:           []*Class{},
	}
}
