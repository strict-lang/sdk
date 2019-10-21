package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type FieldSelectExpression struct {
	Target    StoredExpression
	Selection Expression
	Region    input.Region
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
