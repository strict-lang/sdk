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

type Class struct {
	Name                      string
	Methods                   map[string]*TypedMethod
	Descriptor                TypeDescriptor
	ParameterCount            int
	ImplicitConversionTargets []TypeDescriptor
	SuperType                 *Class
	SuperInterfaces           []*Class
}

// Is returns whether this is the target type or a subtype of the target type.
// In order for this method to return true, one of the following has to be true:
//  - This types descriptor matches the target types descriptor.
//  - This superType is the target type or a subclass of it.
//  - One of the superInterfaces is the target type or a sub interface of it.
//
// If this method returns true, instances of this type can be used as arguments,
// assignments, etc for fields of the target type.
func (class *Class) Is(target Class) bool {
	if class.IsExact(target) || class.SuperType.Is(target) {
		return true
	}
	for _, superInterface := range class.SuperInterfaces {
		if superInterface.Is(target) {
			return true
		}
	}
	return false
}

// IsExact returns whether the target is the exact same class as this.
// It does not return true for supertypes, as opposed to Is(Class).
func (class *Class) IsExact(target Class) bool {
	return class.Descriptor == target.Descriptor
}

// IsImplicitlyConvertibleTo returns whether the class can be implicitly
// converted to the target descriptor.
func (class *Class) IsImplicitlyConvertibleTo(targetDescriptor TypeDescriptor) bool {
	if targetDescriptor == typeDescriptorAny {
		return true
	}
	for _, validTarget := range class.ImplicitConversionTargets {
		if targetDescriptor == validTarget {
			return true
		}
	}
	return false
}

// HasMethodWithName returns true if the class has a method with the passed name.
func (class *Class) HasMethodWithName(name string) bool {
	if _, hasMethod := class.Methods[name]; !hasMethod {
		return class.SuperType.HasMethodWithName(name)
	}
	return true
}

// LookupMethod searches for a method with the passed name.
func (class *Class) LookupMethod(name string) (method *TypedMethod, ok bool) {
	if method, ok = class.Methods[name]; !ok {
		return class.SuperType.LookupMethod(name)
	}
	return
}

func (class *Class) IsParameterized() bool {
	return class.ParameterCount != 0
}
