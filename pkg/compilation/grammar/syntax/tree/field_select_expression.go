package tree

type FieldSelectExpression struct {
	Target       StoredExpression
	Selection    Expression
	NodePosition InputRegion
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

func (expression *FieldSelectExpression) Area() InputRegion {
	return expression.NodePosition
}