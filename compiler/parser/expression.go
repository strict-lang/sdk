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
	next := parser.tokens.Pull()
	return &ast.Identifier{
		Value: next.Value(),
	}, nil
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
		stop, node, err := parser.parseOperationOnOperand(operand)
		if err != nil {
			return operand, err
		}
		if stop {
			break
		}
		operand = node
	}
	return operand, nil
}

// ParseOperationOnOperand parses an operation on an operand that has already
// been parsed. It is called by the ParseOperand method.
func (parser *Parser) parseOperationOnOperand(operand ast.Node) (done bool, node ast.Node, err error) {
	switch next := parser.tokens.Peek(); {
	case token.OperatorValue(next) == token.LeftParenOperator:
		call, err := parser.ParseMethodCall(operand)
		return false, call, err
	default:
		return true, operand, nil
	}
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parser *Parser) ParseBinaryExpression() (ast.BinaryExpression, error) {
	leftOperand, err := parser.ParseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	operator := parser.tokens.Pull()
	if !token.IsOperatorToken(operator) {
		return ast.BinaryExpression{}, ErrInvalidExpression
	}
	rightOperand, err := parser.ParseExpression()
	if err != nil {
		return ast.BinaryExpression{}, err
	}
	return ast.BinaryExpression{
		Operator:     operator.(*token.OperatorToken).Operator,
		LeftOperand:  leftOperand,
		RightOperand: rightOperand,
	}, nil
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
	operand, err := parser.ParseOperation()
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
		Method: method,
	}, nil
}

// parseArgumentList parses the arguments of a MethodCall.
func (parser *Parser) parseArgumentList() ([]ast.Node, error) {
	var arguments []ast.Node
	for {
		next, err := parser.ParseExpression()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, next)
		nextToken := parser.tokens.Pull()
		if !token.IsOperatorToken(nextToken) {
			return arguments, &UnexpectedTokenError{
				Token:    nextToken,
				Expected: "',' or ')'",
			}
		}
		operator := nextToken.(*token.OperatorToken).Operator
		if operator == token.RightParenOperator {
			break
		}
		if operator != token.CommaOperator {
			return arguments, &UnexpectedTokenError{
				Token:    nextToken,
				Expected: "',' or ')'",
			}
		}
	}
	return arguments, nil
}
