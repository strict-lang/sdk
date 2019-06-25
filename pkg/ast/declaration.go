package ast

type Declaration interface {
	declaration()
}

/* Local or private fields of type Type are usually named 'typeNode' if no
   better name is available. This is done because 'type' is a golang keyword. */

// Type represents the type of a member and expressions.
type Type struct {
	position Position
	Name     Identifier
	Members []*Member
}

// Position returns the position of the types declaration.
func (typeNode Type) Position() Position {
	return typeNode.position
}

func (typeNode Type) declaration() {}

type Member struct {
	Typed

	Position  Position
	ValueType *Type
	Name      Identifier
}

func (Member) declaration() {}

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
