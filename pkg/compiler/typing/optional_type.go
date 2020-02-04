package typing

import "fmt"

type OptionalType struct {
	Child Type
}

func (optional *OptionalType) Concrete() Type {
	return optional.Child.Concrete()
}

func (optional *OptionalType) String() string {
	return fmt.Sprintf("%s?", optional.Child)
}

func (optional *OptionalType) Is(target Type) bool {
	if optionalType, ok := target.(*OptionalType); ok {
		return optional.Child.Is(optionalType.Child)
	}
	return false
}

func (optional *OptionalType) Accept(visitor Visitor) {
	visitor.VisitOptional(optional)
}

func (optional *OptionalType) AcceptRecursive(visitor Visitor) {
	optional.Accept(visitor)
	optional.Child.AcceptRecursive(visitor)
}
