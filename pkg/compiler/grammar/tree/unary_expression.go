package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	Operator token.Operator
	Operand  Expression
	Region   input.Region
}

func (unary *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(unary)
}

func (unary *UnaryExpression) AcceptRecursive(visitor Visitor) {
	unary.Accept(visitor)
	unary.Operand.AcceptRecursive(visitor)
}

func (unary *UnaryExpression) Locate() input.Region {
	return unary.Region
}

func (unary *UnaryExpression) Matches(node Node) bool {
	if target, ok := node.(*UnaryExpression); ok {
		return unary.Operator == target.Operator &&
			unary.Operand.Matches(target.Operand)
	}
	return false
}
