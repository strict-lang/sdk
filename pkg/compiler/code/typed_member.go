package code

type MethodDescriptor string

type TypedMethod struct {
	Name       string
	ReturnType *Class
	Parameters []TypedMethod
}

func (method *TypedMethod) Match() {

}

type TypedMethodParameter struct {
	Name string
	Type *Class
}

type TypedInstanceField struct {
}

type TypedLocalField struct {
}
