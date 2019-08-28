package code

const (
	typeDescriptorAny    = "Strict.any"
	typeDescriptorVoid   = "Strict.void"
	typeDescriptorText   = "Strict.text"
	typeDescriptorBool   = "Strict.bool"
	typeDescriptorNumber = "Strict.number"
)

type TypeDescriptor string

var defaultImplicitConversionTargets = []TypeDescriptor{
	typeDescriptorAny, typeDescriptorNumber,
}

type Type struct {
	Name                      string
	Methods                   map[string]TypedMethod
	Descriptor                TypeDescriptor
	ParameterCount            int
	ImplicitConversionTargets []TypeDescriptor
}

type TypeBuilder struct {
	Name                      string
	Descriptor                TypeDescriptor
	ParameterCount            int
	Methods                   []TypedMethod
	ImplicitConversionTargets []TypeDescriptor
}

func (builder TypeBuilder) Create() Type {
	methods := map[string]TypedMethod{}
	for _, method := range builder.Methods {
		methods[method.Name] = method
	}
	return Type{
		Name:           builder.Name,
		Descriptor:     builder.Descriptor,
		ParameterCount: builder.ParameterCount,
		Methods:        methods,
		ImplicitConversionTargets: append(
			builder.ImplicitConversionTargets, defaultImplicitConversionTargets...),
	}
}

func (type_ Type) IsExact(target Type) bool {
	return type_.Descriptor == target.Descriptor
}

func (type_ Type) IsImplicitlyConvertibleTo(targetDescriptor TypeDescriptor) bool {
	if targetDescriptor == typeDescriptorAny {
		return true
	}
	for _, validTarget := range type_.ImplicitConversionTargets {
		if targetDescriptor == validTarget {
			return true
		}
	}
	return false
}

func (type_ Type) HasMethodWithName(name string) (hasMethod bool) {
	_, hasMethod = type_.Methods[name]
	return
}

func (type_ Type) LookupMethod(name string) (method TypedMethod, ok bool) {
	method, ok = type_.Methods[name]
	return
}

func (type_ Type) IsParameterized() bool {
	return type_.ParameterCount != 0
}
