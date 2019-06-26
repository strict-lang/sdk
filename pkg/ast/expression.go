package ast

import "github.com/BenjaminNitschke/Strict/pkg/token"

type Expression interface {
	Node
	expression()
}

// TypedExpression is an Expression whom's return type is known during
// compilation. Examples are arithmetic of numeral constants and
// concatination of text literals.
type TypedExpression interface {
	Typed
	Expression
}

type Identifier struct {
	value    string
	position Position
}

func (identifier Identifier) Position() Position {
	return identifier.position
}

func (identifier Identifier) Accept(visitor Visitor) { }
func (identifier Identifier) expression() { }

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	operator token.Operator
	operand Expression
	position Position
}

func (unary UnaryExpression) Position() Position {
	return unary.position
}

func (unary UnaryExpression) Accept(visitor Visitor) {
	visitor(unary.operand)
}

func (unary UnaryExpression) expression() { }

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	leftOperand Expression
	rightOperand Expression
	operator token.Operator
	position Position
}

func (binary BinaryExpression) Position() Position {
	return binary.position
}

func (binary BinaryExpression) Accept(visitor Visitor) {
	visitor(binary.leftOperand)
	visitor(binary.rightOperand)
}