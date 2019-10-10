package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ListSelectExpression struct {
	Index  Node
	Target Node
	Region input.Region
}

func (expression *ListSelectExpression) Accept(visitor Visitor) {
	VisitListSelectExpression(expression)
}

func (expression *ListSelectExpression) AcceptRecursive(visitor Visitor) {
	expression.Accept(visitor)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (expression *ListSelectExpression) Locate() input.Region {
	return expression.Region
}