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
	VisitConditionalStatement(statement)
}

func (statement *ConditionalStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
	if statement.HasAlternative() {
		AcceptRecursive(visitor)
	}
}

func (statement *ConditionalStatement) Locate() input.Region {
	return statement.Region
}

func (statement *ConditionalStatement) IsModifyingControlFlow() bool {
	return true
}