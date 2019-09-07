package code

type MethodDescriptor string

type TypedMethod struct {
	Name       string
	ReturnType *Type
	Parameters []TypedMethod
}

type TypedMethodParameter struct {
	Name string
	Type *Type
}

type TypedInstanceField struct {
}

type TypedLocalField struct {
}
