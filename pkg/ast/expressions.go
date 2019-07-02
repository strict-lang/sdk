package ast

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

type Identifier struct {
	Value string
}

func NewIdentifier(value string) Identifier {
	return Identifier{Value: value}
}

func (identifier *Identifier) Accept(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier Identifier) String() string {
	return identifier.Value
}

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	Operator token.Operator
	Operand  Node
}

func (unary *UnaryExpression) Accept(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
}

func (unary UnaryExpression) String() string {
	return fmt.Sprintf("UnaryExpression(%s, %s)", unary.Operator, unary.Operand)
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

func (binary BinaryExpression) String() string {
	return fmt.Sprintf("BinaryExpression(%s, %s, %s)", binary.Operator, binary.LeftOperand, binary.RightOperand)
}
