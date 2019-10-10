package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// AssignStatement assigns values to left-hand-side expressions. Operations like
// add-assign are also represented by this Node. If the 'Target' node is a
// FieldDeclaration, this is a field definition.
type AssignStatement struct {
	Target       Node
	Value        Node
	Operator     token.Operator
	Region input.Region
}

func (statement *AssignStatement) Accept(visitor Visitor) {
	visitor.VisitAssignStatement(statement)
}

func (statement *AssignStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	statement.Target.AcceptRecursive(visitor)
	statement.Value.AcceptRecursive(visitor)
}

func (statement *AssignStatement) Locate() input.Region {
	return statement.Region
}