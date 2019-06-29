package ast

import "github.com/BenjaminNitschke/Strict/pkg/token"

type Identifier struct {
	value string
}

func (identifier *Identifier) Accept(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	operator token.Operator
	operand  Node
}

func (unary *UnaryExpression) Accept(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
	unary.operand.Accept(visitor)
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	leftOperand  Node
	rightOperand Node
	operator     token.Operator
}

func (binary *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(binary)
	binary.leftOperand.Accept(visitor)
	binary.rightOperand.Accept(visitor)
}
