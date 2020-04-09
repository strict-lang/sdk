package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() tree.Node {
	parsing.beginStructure(tree.UnknownNodeKind)
	leftHandSide := parsing.parseExpression()
	return parsing.completeInstructionOnNode(leftHandSide)
}

func (parsing *Parsing) completeInstructionOnNode(
	leftHandSide tree.Expression) tree.Node {

	if operator := token.OperatorValue(parsing.token()); operator.IsAssign() {
		return parsing.completeAssignStatement(leftHandSide)
	}
	return parsing.parsePostfixExpression(leftHandSide)
}

func (parsing *Parsing) completeAssignStatement(
	leftHandSide tree.Node) tree.Node {

	parsing.updateTopStructureKind(tree.AssignStatementNodeKind)
	operator := token.OperatorValue(parsing.pullToken())
	rightHandSide := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.AssignStatement{
		Target:   leftHandSide,
		Value:    rightHandSide,
		Operator: operator,
		Region:   parsing.completeStructure(tree.AssignStatementNodeKind),
	}
}

func (parsing *Parsing) parsePostfixExpression(
	leftHandSide tree.Expression) tree.Node {

	if parsing.isLookingAtPostfixExpression() {
		parsing.completePostfixExpressionOnNode(leftHandSide)
	}
	return parsing.completeInvalidPostfixExpression(leftHandSide)
}

func (parsing *Parsing) completeInvalidPostfixExpression(
	leftHandSide tree.Expression) tree.Node {

	parsing.advance()
	parsing.completeStructure(tree.WildcardNodeKind)
	return &tree.ExpressionStatement{
		Expression: leftHandSide,
	}
}

func (parsing *Parsing) isLookingAtPostfixExpression() bool {
	current := parsing.token()
	operator := token.OperatorValue(current)
	return isPostcrementOperator(operator) ||
		token.HasKeywordValue(current, token.ExistsKeyword)
}

func isPostcrementOperator(operator token.Operator) bool {
	return operator == token.IncrementOperator || operator == token.DecrementOperator
}

func (parsing *Parsing) completePostfixExpressionOnNode(
	leftHandSide tree.Expression) tree.Expression {

	operation := token.OperatorValue(parsing.token())
	parsing.advance()
	parsing.skipEndOfStatement()
	return &tree.PostfixExpression{
		Operand:  leftHandSide,
		Operator: operation,
		Region:   parsing.completeStructure(tree.PostfixExpressionNodeKind),
	}
}
