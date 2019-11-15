package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ConditionalStatement struct {
	Condition   Expression
	Alternative Statement
	Consequence Statement
	Region      input.Region
}

func (conditional *ConditionalStatement) HasAlternative() bool {
	return conditional.Alternative != nil
}

func (conditional *ConditionalStatement) Accept(visitor Visitor) {
	visitor.VisitConditionalStatement(conditional)
}

func (conditional *ConditionalStatement) AcceptRecursive(visitor Visitor) {
	conditional.Accept(visitor)
	conditional.Condition.AcceptRecursive(visitor)
	conditional.Consequence.AcceptRecursive(visitor)
	if conditional.HasAlternative() {
		conditional.Alternative.AcceptRecursive(visitor)
	}
}

func (conditional *ConditionalStatement) Locate() input.Region {
	return conditional.Region
}

func (conditional *ConditionalStatement) IsModifyingControlFlow() bool {
	return true
}

func (conditional *ConditionalStatement) Matches(node Node) bool {
	if target, ok := node.(*ConditionalStatement); ok {
		return conditional.matchesStatement(target)
	}
	return false
}

func (conditional *ConditionalStatement) matchesStatement(
	target *ConditionalStatement) bool {

	return conditional.Condition.Matches(target.Condition) &&
		conditional.Consequence.Matches(target.Consequence) &&
		conditional.matchesAlternative(target)
}

func (conditional *ConditionalStatement) matchesAlternative(
	target *ConditionalStatement) bool {

	if !conditional.HasAlternative() {
		return !target.HasAlternative()
	}
	return conditional.Alternative.Matches(target.Alternative)
}
