package ilasm

import "strict.dev/sdk/pkg/compiler/grammar/tree"

func (generation *Generation) generateStatementBlock(
	block *tree.StatementBlock, code *BlockBuilder) {
}

func (generation *Generation) EmitExpressionStatement(
	statement *tree.ExpressionStatement) {

	stack := generation.method.Stack
	previousDepth := stack.CurrentDepth
	generation.EmitExpression(statement.Expression)
	elementsToPop := stack.CurrentDepth - previousDepth
	for count := 0; count < elementsToPop; count++ {
		generation.code.EmitPop()
		stack.DecreaseDepth()
	}
}

func (generation *Generation) EmitBreakStatement(statement *tree.BreakStatement) {
	if label := generation.breakLabel; label != nil {
		generation.code.EmitBranch(label)
	}
}

func (generation *Generation) EmitReturnStatement(
	statement *tree.ReturnStatement) {

	if statement.IsReturningValue() {
		generation.EmitExpression(statement.Value)
		returnClass := resolveClassOfExpression(statement.Value)
		generation.code.EmitValueReturn(returnClass)
	}	else {
		generation.code.EmitReturn()
	}
}

func (generation *Generation) EmitConditionalStatement(
	statement *tree.ConditionalStatement) {

	if statement.HasAlternative() {
		generation.emitConditionalWithAlternative(statement)
	}	else {
		generation.emitConditionalWithoutAlternative(statement)
	}
}

func (generation *Generation) emitConditionalWithAlternative(
	statement *tree.ConditionalStatement) {

	current := generation.code.CreateNextBlock()
	consequence := generation.code.CreateNextBlock()
	alternative :=  consequence.CreateNextBlock()
	exit := alternative.CreateNextBlock()
	generation.EmitExpression(statement.Condition)
	current.EmitBranchIfFalse(exit.Label)
	generation.generateStatementBlock(statement.Consequence, consequence)
	consequence.EmitBranch(exit.Label)
	generation.generateStatementBlock(statement.Alternative, alternative)
	generation.updateCurrentBlock(exit)
}

func (generation *Generation) emitConditionalWithoutAlternative(
	statement *tree.ConditionalStatement) {

	current := generation.code.CreateNextBlock()
	consequence := generation.code.CreateNextBlock()
	exit :=  consequence.CreateNextBlock()
	generation.EmitExpression(statement.Condition)
	current.EmitBranchIfFalse(exit.Label)
	generation.generateStatementBlock(statement.Consequence, consequence)
	generation.updateCurrentBlock(exit)
}