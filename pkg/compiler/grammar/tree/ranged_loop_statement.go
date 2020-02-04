package tree

import "strict.dev/sdk/pkg/compiler/input"

// RangedLoopStatement is a control statement that. Counting from an initial
// value to some target while incrementing a field each step. The values of a
// ranged loop are numeral.
type RangedLoopStatement struct {
	Region input.Region
	Field  *Identifier
	Begin  Expression
	End    Expression
	Body   *StatementBlock
	Parent Node
}

func (loop *RangedLoopStatement) SetEnclosingNode(target Node) {
	loop.Parent = target
}

func (loop *RangedLoopStatement) EnclosingNode() (Node, bool) {
	return loop.Parent, loop.Parent != nil
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

func (loop *RangedLoopStatement) Matches(node Node) bool {
	if target, ok := node.(*RangedLoopStatement); ok {
		return loop.Field.Matches(target.Field) &&
			loop.Begin.Matches(target.Begin) &&
			loop.End.Matches(target.End) &&
			loop.Body.Matches(target.Body)
	}
	return false
}

func (loop *RangedLoopStatement) TransformExpressions(transformer ExpressionTransformer) {
	loop.Begin = loop.Begin.Transform(transformer)
	loop.End = loop.End.Transform(transformer)
}
