package ast

type StringLiteral struct {
	Value string
}

type NumberLiteral struct {
	Value string
}

func (literal *StringLiteral) Accept(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *NumberLiteral) Accept(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}
