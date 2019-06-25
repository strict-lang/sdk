package ast

type Expression interface {
	expression()
}

// TypedExpression is an Expression whom's return type
// is known during compilation. Examples are arithmetic
// of numeral constants and concats of text literals.
type TypedExpression interface {
	Typed
	Expression
}

type Identifier struct {
	Value    string
	Position Position
}

func (id Identifier) expression() {}
