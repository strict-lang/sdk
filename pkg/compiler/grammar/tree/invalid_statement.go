package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// InvalidStatement represents a statement that has not been parsed correctly.
type InvalidStatement struct {
	Region input.Region
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
