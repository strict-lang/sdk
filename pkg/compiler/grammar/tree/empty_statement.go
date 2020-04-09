package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

// EmptyStatement is a statement that does not execute any instructions.
type EmptyStatement struct {
	Region input.Region
	Parent Node
}

func (statement *EmptyStatement) SetEnclosingNode(target Node) {
	statement.Parent = target
}

func (statement *EmptyStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *EmptyStatement) Accept(visitor Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
}

func (statement *EmptyStatement) Locate() input.Region {
	return statement.Region
}

func (statement *EmptyStatement) Matches(node Node) bool {
	_, sameType := node.(*EmptyStatement)
	return sameType
}
