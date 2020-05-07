package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

type ChainExpression struct {
	Expressions []Expression
	Region      input.Region
	Parent      Node
}

func (chain *ChainExpression) LastChild() Expression {
	lastIndex := len(chain.Expressions) - 1
	return chain.Expressions[lastIndex]
}

func (chain *ChainExpression) FirstChild() Expression {
	return chain.Expressions[0]
}

func (chain *ChainExpression) SetEnclosingNode(target Node) {
	chain.Parent = target
}

func (chain *ChainExpression) EnclosingNode() (Node, bool) {
	return chain.Parent, chain.Parent != nil
}

func (chain *ChainExpression) ResolveType(class *scope.Class) {
	chain.LastChild().ResolveType(class)
}

func (chain *ChainExpression) ResolvedType() (*scope.Class, bool) {
	return chain.LastChild().ResolvedType()
}

func (chain *ChainExpression) Accept(visitor Visitor) {
	visitor.VisitFieldSelectExpression(chain)
}

// AcceptRecursive lets the visitor visit the expression and its children.
// The expressions target is accepted prior to the selection.
func (chain *ChainExpression) AcceptRecursive(visitor Visitor) {
	chain.Accept(visitor)
	for _, child := range chain.Expressions {
		child.AcceptRecursive(visitor)
	}
}

func (chain *ChainExpression) Locate() input.Region {
	return chain.Region
}

func (chain *ChainExpression) Matches(node Node) bool {
	if target, ok := node.(*ChainExpression); ok {
		return chain.childrenMatch(target.Expressions)
	}
	return false
}

func (chain *ChainExpression) childrenMatch(target []Expression) bool {
	if len(target) != len(chain.Expressions) {
		return false
	}
	for index, expression := range target {
		if !chain.Expressions[index].Matches(expression) {
			return false
		}
	}
	return true
}

func (chain *ChainExpression) TransformExpressions(transformer ExpressionTransformer) {
	for index, expression := range chain.Expressions {
		chain.Expressions[index] = expression.Transform(transformer)
	}
}

func (chain *ChainExpression) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteFieldSelectExpression(chain)
}
