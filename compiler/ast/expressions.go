package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/token"
)

type Identifier struct {
	Value string
}

func NewIdentifier(value string) *Identifier {
	return &Identifier{Value: value}
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
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token.Operator
}

func (binary *BinaryExpression) Accept(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
}

type SelectorExpression struct {
	Target    Node
	Selection Node
}

func (selector *SelectorExpression) Accept(visitor *Visitor) {
	visitor.VisitSelectorExpression(selector)
}
