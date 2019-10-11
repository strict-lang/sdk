package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// YieldStatement yields an expression to an implicit list that is returned by
// the method it is defined in. Yield statements can only be in methods,
// returning a 'Sequence'. And their values type have to be of the sequences
// element type. Those statements are not accompanied by a ReturnStatement.
type YieldStatement struct {
	Region input.Region
	Value  Expression
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

