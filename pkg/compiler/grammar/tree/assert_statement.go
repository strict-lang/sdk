package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type AssertStatement struct {
	Region     input.Region
	Expression Node
	Parent Node
}

func (assert *AssertStatement) SetEnclosingNode(target Node) {
  assert.Parent = target
}

func (assert *AssertStatement) EnclosingNode() (Node, bool) {
	return assert.Parent, assert.Parent != nil
}

func (assert *AssertStatement) Accept(visitor Visitor) {
	visitor.VisitAssertStatement(assert)
}

func (assert *AssertStatement) AcceptRecursive(visitor Visitor) {
	assert.Accept(visitor)
	assert.Expression.AcceptRecursive(visitor)
}

func (assert *AssertStatement) Locate() input.Region {
	return assert.Region
}

func (assert *AssertStatement) Matches(node Node) bool {
	if target, ok := node.(*AssertStatement); ok {
		return assert.Expression.Matches(target.Expression)
	}
	return false
}
