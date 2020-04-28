package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

type ListExpression struct {
	Expressions []Expression
	Region input.Region
	Parent Node
	resolvedType resolvedType
}

func (list *ListExpression) SetEnclosingNode(target Node) {
	list.Parent = target
}

func (list *ListExpression) EnclosingNode() (Node, bool) {
	return list.Parent, list.Parent != nil
}

func (list *ListExpression) ResolveType(class *scope.Class) {
	list.resolvedType = resolvedType{
		symbol: class,
	}
}

func (list *ListExpression) ResolvedType() (*scope.Class, bool) {
	return list.resolvedType.class()
}

func (list *ListExpression) Accept(visitor Visitor) {
	visitor.VisitListExpression(list)
}

func (list *ListExpression) AcceptRecursive(visitor Visitor) {
	list.Accept(visitor)
	for _, child := range list.Expressions {
		child.AcceptRecursive(visitor)
	}
}

func (list *ListExpression) Locate() input.Region {
	return list.Region
}

func (list *ListExpression) Matches(node Node) bool {
	if target, ok := node.(*ListExpression); ok {
		return list.childrenMatch(target.Expressions)
	}
	return false
}

func (list *ListExpression) childrenMatch(target []Expression) bool {
	if len(target) != len(list.Expressions) {
		return false
	}
	for index, expression := range target {
		if !list.Expressions[index].Matches(expression) {
			return false
		}
	}
	return true
}

func (list *ListExpression) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteListExpression(list)
}