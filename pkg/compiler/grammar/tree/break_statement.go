package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type BreakStatement struct {
	Region     input.Region
	Expression Node
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