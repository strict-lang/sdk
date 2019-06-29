package ast

// Type represents the type of a member and expressions.
type Type struct {
	position Position
	Name     Identifier
	Members  []*Member
}

// Member is a typed field of a class. It represents methods and
// attributes. The type of a method member is its return-type.
type Member struct {
	Name      Identifier
	Position  Position
	ValueType *Type
}

// Position returns the position of the types declaration.
func (typeNode Type) Position() Position {
	return typeNode.position
}

func (member Member) Kind() Kind {
	return TypeNode
}

func (member Member) Type() *Type {
	return member.ValueType
}

type Parameter struct {
	Position Position
	Type     Type
	Named    Identifier
}

type Method struct {
	Member
	Parameters []Parameter
}
