package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// parseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parsing *Parsing) parseBinaryExpression(
	precedence token.Precedence) tree.Expression {
	parsing.beginStructure(tree.BinaryExpressionNodeKind)
	leftHandSide := parsing.parseUnaryExpression()
	for !parsing.isAtEndOfBinaryExpression(precedence) {
		leftHandSide = parsing.parseBinaryExpressionWithLeftHandSide(leftHandSide, precedence)
	}
	parsing.completeStructure(tree.BinaryExpressionNodeKind)
	return leftHandSide
}

func (parsing *Parsing) isAtEndOfBinaryExpression(precedence token.Precedence) bool {
	return token.OperatorValue(parsing.token()).Precedence() < precedence
}

func (parsing *Parsing) parseBinaryExpressionWithLeftHandSide(
	leftHandSide tree.Expression,
	precedence token.Precedence) tree.Expression {

	operator := token.OperatorValue(parsing.token())
	if operator.Precedence() < precedence {
		return leftHandSide
	}
	parsing.advance()
	nextPrecedence := operator.Precedence().Next()
	rightHandSide := parsing.parseBinaryExpression(nextPrecedence)
	return &tree.BinaryExpression{
		Operator:     operator,
		LeftOperand:  leftHandSide,
		RightOperand: rightHandSide,
		Region: input.CreateRegion(
			leftHandSide.Locate().Begin(),
			rightHandSide.Locate().End()),
	}
}
