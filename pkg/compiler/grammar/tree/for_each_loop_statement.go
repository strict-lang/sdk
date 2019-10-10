package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// ForEachLoopStatement is a control statement. Iterating an enumeration without
// requiring explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type that implements the 'Sequence' interface.
type ForEachLoopStatement struct {
	Region input.Region
	Body         Statement
	Sequence     Expression
	Field        *Identifier
}

func (loop *ForEachLoopStatement) Accept(visitor Visitor) {
	visitor.VisitForEachLoopStatement(loop)
}

func (loop *ForEachLoopStatement) AcceptRecursive(visitor Visitor) {
	loop.Accept(visitor)
	loop.Field.AcceptRecursive(visitor)
	loop.Sequence.AcceptRecursive(visitor)
	loop.Body.AcceptRecursive(visitor)
}

func (loop *ForEachLoopStatement) Locate() input.Region {
	return loop.Region
}
