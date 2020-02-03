package tree

import "strict.dev/sdk/pkg/compiler/input"

type BreakStatement struct {
	Region     input.Region
	Parent     Node
}

func (statement *BreakStatement) SetEnclosingNode(target Node) {
  statement.Parent = target
}

func (statement *BreakStatement) EnclosingNode() (Node, bool) {
  return statement.Parent, statement.Parent != nil
}

func (statement *BreakStatement) Accept(visitor Visitor) {
	visitor.VisitBreakStatement(statement)
}

func (statement *BreakStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
}

func (statement *BreakStatement) Locate() input.Region {
	return statement.Region
}

func (statement *BreakStatement) Matches(node Node) bool {
	_, ok := node.(*BreakStatement)
	return ok
}