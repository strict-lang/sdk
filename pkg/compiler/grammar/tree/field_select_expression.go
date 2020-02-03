package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/scope"
)

type FieldSelectExpression struct {
	Target    StoredExpression
	Selection Expression
	Region    input.Region
	Parent Node
	resolvedType resolvedType
}

func (expression *FieldSelectExpression) SetEnclosingNode(target Node) {
  expression.Parent = target
}

func (expression *FieldSelectExpression) EnclosingNode() (Node, bool) {
  return expression.Parent, expression.Parent != nil
}

func (expression *FieldSelectExpression) ResolveType(class *scope.Class) {
  expression.resolvedType.resolve(class)
}

func (expression *FieldSelectExpression) ResolvedType() (*scope.Class, bool) {
  return expression.resolvedType.class()
}

func (expression *FieldSelectExpression) Accept(visitor Visitor) {
	visitor.VisitFieldSelectExpression(expression)
}

// AcceptRecursive lets the visitor visit the expression and its children.
// The expressions target is accepted prior to the selection.
func (expression *FieldSelectExpression) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	expression.Target.AcceptRecursive(visitor)
	expression.Selection.AcceptRecursive(visitor)
}

func (expression *FieldSelectExpression) Locate() input.Region {
	return expression.Region
}

func (expression *FieldSelectExpression) Matches(node Node) bool {
	if target, ok := node.(*FieldSelectExpression); ok {
		return expression.Target.Matches(target.Target) &&
			expression.Selection.Matches(target.Selection)
	}
	return false
}

func (expression *FieldSelectExpression) FindLastIdentifier() (*Identifier, bool) {
	switch expression.Selection.(type) {
		case *Identifier:
			identifier, ok := expression.Selection.(*Identifier)
			return identifier, ok
		case *FieldSelectExpression:
			if next, ok := expression.Selection.(*FieldSelectExpression); ok {
				return next.FindLastIdentifier()
			}
	}
	return nil, false
}

func (expression *FieldSelectExpression) TransformExpressions(
	transformer ExpressionTransformer) {

	expression.Target = expression.Target.Transform(transformer)
	expression.Selection = expression.Selection.Transform(transformer)
}

func (expression *FieldSelectExpression) Transform(
	transformer ExpressionTransformer) Expression {

	return transformer.RewriteFieldSelectExpression(expression)
}