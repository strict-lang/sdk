package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ConditionalStatement struct {
	Condition   Expression
	Alternative Statement
	Consequence Statement
	Region      input.Region
}

func (statement *ConditionalStatement) HasAlternative() bool {
	return statement.Alternative != nil
}

func (statement *ConditionalStatement) Accept(visitor Visitor) {
	visitor.VisitConditionalStatement(statement)
}

func (statement *ConditionalStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	statement.Condition.AcceptRecursive(visitor)
	statement.Consequence.AcceptRecursive(visitor)
	if statement.HasAlternative() {
		statement.Alternative.AcceptRecursive(visitor)
	}
}

func (statement *ConditionalStatement) Locate() input.Region {
	return statement.Region
}

func (statement *ConditionalStatement) IsModifyingControlFlow() bool {
	return true
}
