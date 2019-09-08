package ast

type StringLiteral struct {
	Value        string
	NodePosition Position
}

func (literal *StringLiteral) Accept(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *StringLiteral) AcceptRecursive(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *StringLiteral) Position() Position {
	return literal.NodePosition
}

type NumberLiteral struct {
	Value        string
	NodePosition Position
}

func (literal *NumberLiteral) IsFloat() bool {
	return false
}

func (literal *NumberLiteral) Accept(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}

func (literal *NumberLiteral) AcceptRecursive(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}

func (literal *NumberLiteral) Position() Position {
	return literal.NodePosition
}
