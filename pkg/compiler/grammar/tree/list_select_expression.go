package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ListSelectExpression struct {
	Index  Node
	Target Node
	Region input.Region
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