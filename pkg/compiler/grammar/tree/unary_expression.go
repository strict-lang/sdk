package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	Operator token.Operator
	Operand  Expression
	Region   input.Region
	Parent   Node
	resolvedType resolvedType
}

func (unary *UnaryExpression) ResolveType(class *scope.Class) {
  unary.resolvedType.resolve(class)
}

func (unary *UnaryExpression) ResolvedType() (*scope.Class, bool) {
  return unary.resolvedType.class()
}

func (unary *UnaryExpression) SetEnclosingNode(target Node) {
  unary.Parent = target
}

func (unary *UnaryExpression) EnclosingNode() (Node, bool) {
  return unary.Parent, unary.Parent != nil
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
