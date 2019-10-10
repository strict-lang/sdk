package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

type ListSelectExpression struct {
	Index        Node
	Target       Node
	Region input.Region
}

func (expression *ListSelectExpression) Accept(visitor Visitor) {
	visitor.VisitListSelectExpression(expression)
}

func (expression *ListSelectExpression) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	expression.Index.AcceptRecursive(visitor)
	expression.Target.AcceptRecursive(visitor)
}

func (expression *ListSelectExpression) Locate() input.Region {
	return expression.Region
}