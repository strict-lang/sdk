package tree

import "strict.dev/sdk/pkg/compiler/input"

// InvalidStatement represents a statement that has not been parsed correctly.
type InvalidStatement struct {
	Region input.Region
	Parent Node
}

func (statement *InvalidStatement) SetEnclosingNode(target Node) {
	statement.Parent = target
}

func (statement *InvalidStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *InvalidStatement) Accept(visitor Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
}

func (statement *InvalidStatement) Locate() input.Region {
	return statement.Region
}

func (statement *InvalidStatement) Matches(node Node) bool {
	_, sameType := node.(*InvalidStatement)
	return sameType
}
