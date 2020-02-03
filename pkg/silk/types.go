package silk

import "strict.dev/sdk/pkg/silk/symbol"

type Type interface {
	ClassReference() symbol.Reference
	Matches(Type) bool
}

type SliceReference struct {
	Class symbol.Reference
}

func (reference *SliceReference) ClassReference() symbol.Reference {
	return reference.Class
}

func CreateSliceReferenceType(elementClass symbol.Reference) *SliceReference {
	return &SliceReference{
		Class: elementClass,
	}
}

type Reference struct {
	Class symbol.Reference
}

func (reference *Reference) ClassReference() symbol.Reference {
	return reference.Class
}

func CreateReferenceType(class symbol.Reference) *Reference {
	return &Reference{
		Class: class,
	}
}

type Primitive struct {
	Class symbol.Reference
}

func (primitive Primitive) ClassReference() symbol.Reference {
	return primitive.Class
}

func (primitive Primitive) Matches(target Type) bool {
	return false
}

var VoidType = Primitive{Class: 0}
