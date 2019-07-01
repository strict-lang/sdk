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
	Operator token.Operator
	Operand  Node
}

func (unary *UnaryExpression) Accept(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
	unary.Operand.Accept(visitor)
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token.Operator
}

func (binary *BinaryExpression) Accept(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
	binary.LeftOperand.Accept(visitor)
	binary.RightOperand.Accept(visitor)
}
