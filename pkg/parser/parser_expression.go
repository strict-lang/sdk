package parser

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/token"

	"github.com/BenjaminNitschke/Strict/pkg/ast"
)

var (
	// ErrInvalidExpression is returned from a function that fails to parse
	// an expression. Functions returning this should report more verbose
	// error messages to the diagnostics.Recorder.
	ErrInvalidExpression = errors.New("could not parse invalid expression")
)

func (parser Parser) ParseExpression() (ast.Node, error) {
	return nil, nil
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parser Parser) ParseBinaryExpression() (ast.BinaryExpression, error) {
	leftOperand, err := parser.ParseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	operator := parser.tokens.Peek()
	if !operator.IsOperator() {
		return ast.BinaryExpression{}, ErrInvalidExpression
	}
	rightOperand, err := parser.ParseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	return ast.BinaryExpression{
		Operator:     operator.(token.OperatorToken).Operator,
		LeftOperand:  leftOperand,
		RightOperand: rightOperand,
	}, nil
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parser Parser) ParseUnaryExpression() (ast.UnaryExpression, error) {
	operator := parser.tokens.Pull()
	if !operator.IsOperator() {
		return ast.UnaryExpression{}, ErrInvalidExpression
	}
	expression, err := parser.ParseExpression()
	if err != nil {
		return ast.UnaryExpression{}, err
	}
	return ast.UnaryExpression{
		Operator: operator.(token.OperatorToken).Operator,
		Operand:  expression,
	}, nil
}
