package ast

type Literal interface {
	Expression
	literal()
}

type StringLiteral struct {
	value string
	position Position
}

func (literal StringLiteral) literal() {}
func (literal StringLiteral) expression() {}

func (literal StringLiteral) Position() Position {
	return literal.position
}

type NumberLiteral struct {
	value string
	position Position
}

func (literal NumberLiteral) literal() {}
func (literal NumberLiteral) expression() {}

func (literal NumberLiteral) Position() Position {
	return literal.position
}
