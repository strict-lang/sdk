package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

// YieldStatement yields an expression to an implicit list that is returned by
// the method it is defined in. Yield statements can only be in methods,
// returning a 'Sequence'. And their values type have to be of the sequences
// element type. Those statements are not accompanied by a ReturnStatement.
type YieldStatement struct {
	Region input.Region
	Value  Expression
	Parent Node
}

func (yield *YieldStatement) SetEnclosingNode(target Node) {
	yield.Parent = target
}

func (yield *YieldStatement) EnclosingNode() (Node, bool) {
	return yield.Parent, yield.Parent != nil
}

func (yield *YieldStatement) Accept(visitor Visitor) {
	visitor.VisitYieldStatement(yield)
}

func (yield *YieldStatement) AcceptRecursive(visitor Visitor) {
	yield.Accept(visitor)
	yield.Value.AcceptRecursive(visitor)
}

func (yield *YieldStatement) Locate() input.Region {
	return yield.Region
}

func (yield *YieldStatement) Matches(node Node) bool {
	if target, ok := node.(*YieldStatement); ok {
		return yield.Value.Matches(target.Value)
	}
	return false
}

func (yield *YieldStatement) TransformExpressions(transformer ExpressionTransformer) {
	yield.Value = yield.Value.Transform(transformer)
}
