package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	expression.Expression.AcceptRecursive(visitor)
}

func (expression *ExpressionStatement) Locate() input.Region {
	return expression.Expression.Locate()
}

func (expression *ExpressionStatement) Matches(node Node) bool {
	if target, ok := node.(*ExpressionStatement); ok {
		return expression.Expression.Matches(target.Expression)
	}
	return false
}
