package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// RangedLoopStatement is a control statement that. Counting from an initial
// value to some target while incrementing a field each step. The values of a
// ranged loop are numeral.
type RangedLoopStatement struct {
	Region input.Region
	Field   *Identifier
	Begin Expression
	End     Expression
	Body         Statement
}

func (loop *RangedLoopStatement) Accept(visitor Visitor) {
	visitor.VisitRangedLoopStatement(loop)
}

func (loop *RangedLoopStatement) AcceptRecursive(visitor Visitor) {
	loop.Accept(visitor)
	loop.Field.AcceptRecursive(visitor)
	loop.Begin.AcceptRecursive(visitor)
	loop.End.AcceptRecursive(visitor)
	loop.Body.AcceptRecursive(visitor)
}

func (loop *RangedLoopStatement) Locate() input.Region {
	return loop.Region
}

