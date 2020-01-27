package typing

type Type interface {
	Is(Type) bool
	Accept(Visitor)
	AcceptRecursive(Visitor)
	Concrete() Type
	String() string
}

type Visitor interface {
	VisitList(*ListType)
	VisitGeneric(*GenericType)
	VisitConcrete(*ConcreteType)
	VisitOptional(*OptionalType)
}

func NewEmptyClass(name string) Type {
	return nil
}