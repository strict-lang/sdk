package ast

// Type represents the type of a member and expressions.
type Type struct {
	Name    *Identifier
	Members []*Member
}

func (typeNode *Type) Accept(visitor *Visitor) {
	visitor.VisitType(typeNode)
}

// Member is a typed field of a class. It represents methods and
// attributes. The type of a method member is its return-type.
type Member struct {
	Name   Identifier
	Type  TypeName
}

func (member *Member) Accept(visitor *Visitor) {
	visitor.VisitMember(member)
}

type Method struct {
	Name Identifier
	Type TypeName
	Parameters []Parameter
	Body Node
}

func (method *Method) Accept(visitor *Visitor) {
	visitor.VisitMethod(method)
}

type Parameter struct {
	Type TypeName
	Name Identifier
}

func (parameter *Parameter) Accept(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}
