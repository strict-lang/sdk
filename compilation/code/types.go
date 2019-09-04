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
	SuperType                 *Type
	SuperInterfaces           []Type
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

// Is returns whether this is the target type or a subtype of the target type.
// In order for this method to return true, one of the following has to be true:
//  - This types descriptor matches the target types descriptor.
//  - This superType is the target type or a subclass of it.
//  - One of the superInterfaces is the target type or a sub interface of it.
//
// If this method returns true, instances of this type can be used as arguments,
// assignments, etc for fields of the target type.
func (type_ Type) Is(target Type) bool {
	if type_.IsExact(target) || type_.SuperType.Is(target) {
		return true
	}
	for _, superInterface := range type_.SuperInterfaces {
		if superInterface.Is(target) {
			return true
		}
	}
	return false
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

func (type_ Type) HasMethodWithName(name string) bool {
	if _, hasMethod := type_.Methods[name]; !hasMethod {
		return type_.SuperType.HasMethodWithName(name)
	}
	return true
}

func (type_ Type) LookupMethod(name string) (method TypedMethod, ok bool) {
	if method, ok = type_.Methods[name]; !ok {
		return type_.SuperType.LookupMethod(name)
	}
	return
}

func (type_ Type) IsParameterized() bool {
	return type_.ParameterCount != 0
}
