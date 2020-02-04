package tree

import "strict.dev/sdk/pkg/compiler/input"

// ForEachLoopStatement is a control statement. Iterating an enumeration without
// requiring explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type that implements the 'Sequence' interface.
type ForEachLoopStatement struct {
	Region   input.Region
	Body     *StatementBlock
	Sequence Expression
	Field    *Identifier
	Parent   Node
}

func (loop *ForEachLoopStatement) SetEnclosingNode(target Node) {
	loop.Parent = target
}

func (loop *ForEachLoopStatement) EnclosingNode() (Node, bool) {
	return loop.Parent, loop.Parent != nil
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

func (loop *ForEachLoopStatement) Matches(node Node) bool {
	if target, ok := node.(*ForEachLoopStatement); ok {
		return loop.Field.Matches(target.Field) &&
			loop.Sequence.Matches(target.Sequence) &&
			loop.Body.Matches(target.Body)
	}
	return false
}

func (loop *ForEachLoopStatement) TransformExpressions(transformer ExpressionTransformer) {
	loop.Sequence = loop.Sequence.Transform(transformer)
}
