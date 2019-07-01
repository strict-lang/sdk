package ast

// Type represents the type of a member and expressions.
type Type struct {
	Name    *Identifier
	Members []*Member
}

func (typeNode *Type) Accept(visitor *Visitor) {
	visitor.VisitType(typeNode)
	for _, member := range typeNode.Members {
		member.Accept(visitor)
	}
}

// Member is a typed field of a class. It represents methods and
// attributes. The type of a method member is its return-type.
type Member struct {
	Name      Identifier
	ValueType *Type
}

func (member Member) Type() *Type {
	return member.ValueType
}

func (member *Member) Accept(visitor *Visitor) {
	visitor.VisitMember(member)
}

type Method struct {
	Member
	Parameters []Parameter
}

func (method *Method) Accept(visitor *Visitor) {
	visitor.VisitMethod(method)
	for _, parameter := range method.Parameters {
		parameter.Accept(visitor)
	}
}

type Parameter struct {
	Type     Type
	Named    Identifier
}

func (parameter *Parameter) Accept(visitor *Visitor) {
	visitor.VisitParameter(parameter)
}
