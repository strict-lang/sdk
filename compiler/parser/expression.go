package parser

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/compiler/token"

	"github.com/BenjaminNitschke/Strict/compiler/ast"
)

var (
	// ErrInvalidExpression is returned from a function that fails to parse
	// an expression. Functions returning this should report more verbose
	// error messages to the diagnostics.Recorder.
	ErrInvalidExpression = errors.New("could not parse invalid expression")
)

func (parser *Parser) ParseExpression() (ast.Node, error) {
	parser.tokens.Pull()
	return parser.parseBinaryExpression(token.LowPrecedence + 1)
}

func (parser *Parser) ParseOperand() (ast.Node, error) {
	switch last := parser.tokens.Last(); {
	case token.IsIdentifierToken(last):
		return &ast.Identifier{Value: last.Value()}, nil
	case token.IsStringLiteralToken(last):
		return &ast.StringLiteral{Value: last.Value()}, nil
	case token.IsNumberLiteralToken(last):
		return &ast.NumberLiteral{Value: last.Value()}, nil
	case token.OperatorValue(last) == token.LeftParenOperator:
		return parser.completeLeftParenExpression()
	}
	return nil, ErrInvalidExpression
}

func (parser *Parser) completeLeftParenExpression() (ast.Node, error) {
	parser.expressionDepth++
	expression, err := parser.ParseExpression()
	if err != nil {
		return expression, err
	}
	parser.expressionDepth--
	if token.OperatorValue(parser.tokens.Last()) != token.RightParenOperator {
		return nil, &UnexpectedTokenError{
			Token:    parser.tokens.Last(),
			Expected: token.RightParenOperator.String(),
		}
	}
	return expression, nil
}

// ParseOperation parses the initial operand and continues to parse operands on
// that operand, forming a node for another expression.
func (parser *Parser) ParseOperation() (ast.Node, error) {
	operand, err := parser.ParseOperand()
	if err != nil {
		return nil, err
	}
	// TODO(merlinosayimwen): Add field selector
	for {
		done, node, err := parser.parseOperationOnOperand(operand)
		if err != nil {
			return operand, err
		}
		operand = node
		if done {
			return operand, nil
		}
	}
}

// ParseOperationOnOperand parses an operation on an operand that has already
// been parsed. It is called by the ParseOperand method.
func (parser *Parser) parseOperationOnOperand(operand ast.Node) (done bool, node ast.Node, err error) {
	switch next := parser.tokens.Peek(); {
	case token.OperatorValue(next) == token.LeftParenOperator:
		call, err := parser.ParseMethodCall(operand)
		return false, call, err
	}
	return true, operand, nil
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parser *Parser) parseBinaryExpression(requiredPrecedence token.Precedence) (ast.Node, error) {
	leftHandSide, err := parser.ParseUnaryExpression()
	if err != nil {
		return nil, err
	}
	for {
		next := parser.tokens.Pull()
		precedence := token.PrecedenceOfAny(next)
		if precedence < requiredPrecedence {
			return leftHandSide, nil
		}
		parser.tokens.Pull()
		rightHandSide, err := parser.parseBinaryExpression(precedence)
		if err != nil {
			return leftHandSide, err
		}
		leftHandSide = &ast.BinaryExpression{
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			Operator:     token.OperatorValue(next),
		}
	}
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parser *Parser) ParseUnaryExpression() (ast.Node, error) {
	operatorToken := parser.tokens.Last()
	if !token.IsOperatorOrOperatorKeywordToken(operatorToken) {
		return parser.ParseOperation()
	}
	operator := token.OperatorValue(operatorToken)
	if !operator.IsUnaryOperator() {
		return parser.ParseOperation()
	}
	parser.tokens.Pull()
	operand, err := parser.ParseUnaryExpression()
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}, nil
}

// ParseMethodCall parses the call to a method.
func (parser *Parser) ParseMethodCall(method ast.Node) (*ast.MethodCall, error) {
	if err := parser.skipOperator(token.LeftParenOperator); err != nil {
		return &ast.MethodCall{}, err
	}
	arguments, err := parser.parseArgumentList()
	if err != nil {
		return &ast.MethodCall{}, err
	}
	return &ast.MethodCall{
		Arguments: arguments,
		Method:    method,
	}, nil
}

// parseArgumentList parses the arguments of a MethodCall.
func (parser *Parser) parseArgumentList() ([]ast.Node, error) {
	if token.OperatorValue(parser.tokens.Peek()) == token.RightParenOperator {
		parser.tokens.Pull()
		return []ast.Node{}, nil
	}
	var arguments []ast.Node
	for {
		next, err := parser.ParseExpression()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, next)
		nextToken := parser.tokens.Last()
		switch token.OperatorValue(nextToken) {
		case token.RightParenOperator:
			return arguments, nil
		case token.CommaOperator:
			continue
		}
		return arguments, &UnexpectedTokenError{
			Token:    nextToken,
			Expected: "',' or ')'",
		}
	}
}
