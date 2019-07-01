package ast

type StringLiteral struct {
	value string
}

type NumberLiteral struct {
	value string
}

func (literal *StringLiteral) Accept(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *NumberLiteral) Accept(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}
