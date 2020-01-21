package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

type ListSelectExpression struct {
	Index        Expression
	Target       Expression
	Region       input.Region
	resolvedType resolvedType
	Parent Node
}

func (expression *ListSelectExpression) SetEnclosingNode(target Node) {
  expression.Parent = target
}

func (expression *ListSelectExpression) EnclosingNode() (Node, bool) {
  return expression.Parent, expression.Parent != nil
}

func (expression *ListSelectExpression) ResolveType(class *scope.Class) {
  expression.resolvedType.resolve(class)
}

func (expression *ListSelectExpression) ResolvedType() (*scope.Class, bool) {
  return expression.resolvedType.class()
}

func (expression *ListSelectExpression) Accept(visitor Visitor) {
	visitor.VisitListSelectExpression(expression)
}

func (expression *ListSelectExpression) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	expression.Target.AcceptRecursive(visitor)
	expression.Index.AcceptRecursive(visitor)
}

func (expression *ListSelectExpression) Locate() input.Region {
	return expression.Region
}

func (expression *ListSelectExpression) Matches(node Node) bool {
	if target, ok := node.(*ListSelectExpression); ok {
		return expression.Index.Matches(target.Index) &&
			expression.Target.Matches(target.Target)
	}
	return false
}

func (expression *ListSelectExpression) TransformExpressions(
	transformer ExpressionTransformer) {

	expression.Index = expression.Index.Transform(transformer)
	expression.Target = expression.Target.Transform(transformer)
}

func (expression *ListSelectExpression) Transform(
	transformer ExpressionTransformer) Expression {

	return transformer.RewriteListSelectExpression(expression)
}