package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ListSelectExpression struct {
	Index  Node
	Target Node
	Region input.Region
	resolvedType resolvedType
}

func (expression *ListSelectExpression) Resolve(descriptor TypeDescriptor) {
	expression.resolvedType.setDescriptor(descriptor)
}

func (expression *ListSelectExpression) GetResolvedType() (TypeDescriptor, bool) {
	return expression.resolvedType.getDescriptor()
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
