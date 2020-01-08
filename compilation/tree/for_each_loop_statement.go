package tree

import "gitlab.com/strict-lang/sdk/compilation/input"

// ForEachLoopStatement is a control statement. Iterating an enumeration without
// requiring explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type that implements the 'Sequence' interface.
type ForEachLoopStatement struct {
	Region   input.Region
	Body     Statement
	Sequence Expression
	Field    *Identifier
}

func (loop *ForEachLoopStatement) Accept(visitor Visitor) {
	VisitForEachLoopStatement(loop)
}

func (loop *ForEachLoopStatement) AcceptRecursive(visitor Visitor) {
	loop.Accept(visitor)
	loop.Field.AcceptRecursive(visitor)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (loop *ForEachLoopStatement) Locate() input.Region {
	return loop.Region
}

func (loop *ForEachLoopStatement) Matches(node Node) bool {
	if target, ok := node.(*ForEachLoopStatement); ok {
		return loop.Field.Matches(target.Field) &&
			Matches(target.Sequence) &&
			Matches(target.Body)
	}
	return false
}