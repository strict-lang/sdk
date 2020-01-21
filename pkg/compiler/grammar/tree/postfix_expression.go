package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

// PostfixExpression is an expression with an operator that is written in a
// reverse polish notation. The only postfix operations in strict are a
// PostIncrement and PostDecrement. The supported postfix expressions modify
// the value of their operand, thus their operand can not be of a constant value.
type PostfixExpression struct {
	// Operand is the expression that is modified by this expression. It can not
	// be immutable and has to be stored.
	Operand StoredExpression
	// Operator is the type of operation that is applied to the operand.
	Operator token.Operator
	// InputRegion is the area of code covered by the node.
	Region       input.Region
	resolvedType resolvedType
	Parent Node
}

func (expression *PostfixExpression) SetEnclosingNode(target Node) {
  expression.Parent = target
}

func (expression *PostfixExpression) EnclosingNode() (Node, bool) {
  return expression.Parent, expression.Parent != nil
}

func (expression *PostfixExpression) ResolveType(class *scope.Class) {
	expression.resolvedType.resolve(class)
}

func (expression *PostfixExpression) ResolvedType() (*scope.Class, bool) {
	return expression.resolvedType.class()
}

// Accept lets the visitor visit this expression.
func (expression *PostfixExpression) Accept(visitor Visitor) {
	visitor.VisitPostfixExpression(expression)
}

// AcceptRecursive lets the visitor visit the expression and calls the
// same method on every child. Thus the complete branch is visited.
func (expression *PostfixExpression) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	expression.Operand.AcceptRecursive(visitor)
}

// Locate returns the area of code that is covered by the node.
func (expression *PostfixExpression) Locate() input.Region {
	return expression.Region
}

func (expression *PostfixExpression) Matches(node Node) bool {
	if target, ok := node.(*PostfixExpression); ok {
		return expression.Operator == target.Operator &&
			target.Operand.Matches(target.Operand)
	}
	return false
}
