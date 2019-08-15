package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/token"
)

type Identifier struct {
	Value string
	NodePosition Position
}

func (identifier *Identifier) Accept(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) AcceptAll(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) Position() Position {
	return identifier.Position()
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
