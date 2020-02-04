package typing

type ConcreteType struct {
	Name   string
	Traits []Type
}

func (concrete *ConcreteType) Concrete() Type {
	return concrete
}

func (concrete *ConcreteType) String() string {
	return concrete.Name
}

func (concrete *ConcreteType) Is(target Type) bool {
	if targetConcrete, ok := target.(*ConcreteType); ok {
		return concrete.matches(targetConcrete)
	}
	return false
}

func (concrete *ConcreteType) matches(target *ConcreteType) bool {
	if concrete.matchesExact(target) {
		return true
	}
	// FIXME: Use hashtable instead for performance improvement.
	for _, targetTrait := range target.Traits {
		for _, trait := range concrete.Traits {
			if trait.Is(targetTrait) {
				return true
			}
		}
	}
	return false
}

func (concrete *ConcreteType) matchesExact(target *ConcreteType) bool {
	return concrete.Name == target.Name
}

func (concrete *ConcreteType) Accept(visitor Visitor) {
	visitor.VisitConcrete(concrete)
}

func (concrete *ConcreteType) AcceptRecursive(visitor Visitor) {
	concrete.Accept(visitor)
}
