// The expression file contains methods that are parsing expressions. Every method expects
// the first token that it requires to be the current one (parser.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the parser.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

func (parser *Parser) ParseExpression() (ast.Node, error) {
	return parser.parseBinaryExpression(token.LowPrecedence + 1)
}

func (parser *Parser) ParseOperand() (ast.Node, error) {
	switch last := parser.token(); {
	case token.IsIdentifierToken(last): return parser.parseIdentifier()
	case token.IsStringLiteralToken(last): return parser.parseStringLiteral()
	case token.IsNumberLiteralToken(last): return parser.parseNumberLiteral()
	case token.OperatorValue(last) == token.LeftParenOperator:
		return parser.completeLeftParenExpression()
	}
	return nil, ErrInvalidExpression
}

func (parser *Parser) parseIdentifier() (*ast.Identifier, error) {
	defer parser.advance()
	return &ast.Identifier{
		Value: parser.token().Value(),
		NodePosition: parser.createTokenPosition(),
	}, nil
}

func (parser *Parser) parseStringLiteral() (*ast.StringLiteral, error) {
	defer parser.advance()
	return &ast.StringLiteral{
		Value: parser.token().Value(),
		NodePosition: parser.createTokenPosition(),
	}, nil
}

func (parser *Parser) parseNumberLiteral() (*ast.NumberLiteral, error) {
	defer parser.advance()
	return &ast.NumberLiteral{
		Value: parser.token().Value(),
		NodePosition: parser.createTokenPosition(),
	}, nil
}

func (parser *Parser) completeLeftParenExpression() (ast.Node, error) {
	parser.advance()
	parser.expressionDepth++
	expression, err := parser.ParseExpression()
	if err != nil {
		return expression, err
	}
	parser.expressionDepth--
	if token.OperatorValue(parser.token()) != token.RightParenOperator {
		return nil, &UnexpectedTokenError{
			Token:    parser.token(),
			Expected: token.RightParenOperator.String(),
		}
	}
	parser.advance()
	return expression, nil
}

// ParseOperation parses the initial operand and continues to parse operands on
// that operand, forming a node for another expression.
func (parser *Parser) ParseOperation() (ast.Node, error) {
	operand, err := parser.ParseOperand()
	if err != nil {
		return nil, err
	}
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
	switch next := parser.token(); {
	case token.OperatorValue(next) == token.LeftParenOperator:
		call, err := parser.parseMethodCallOnNode(operand)
		return false, call, err
	case token.OperatorValue(next) == token.DotOperator:
		selector, err := parser.parseSelection(operand)
		return false, selector, err
	}
	return true, operand, nil
}

func (parser *Parser) parseSelection(operand ast.Node) (ast.Node, error) {
	beginOffset := parser.offset()
	if err := parser.skipOperator(token.DotOperator); err != nil {
		return nil, err
	}
	field, err := parser.ParseOperand()
	if err != nil {
		return nil, err
	}
	return &ast.SelectorExpression{
		Target:    operand,
		Selection: field,
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parser *Parser) parseBinaryExpression(requiredPrecedence token.Precedence) (ast.Node, error) {
	beginOffset := parser.offset()
	leftHandSide, err := parser.ParseUnaryExpression()
	if err != nil {
		return nil, err
	}
	for {
		operator := parser.token()
		precedence := token.PrecedenceOfAny(operator)
		if precedence < requiredPrecedence {
			return leftHandSide, nil
		}
		parser.advance()
		rightHandSide, err := parser.parseBinaryExpression(precedence)
		if err != nil {
			return leftHandSide, err
		}
		leftHandSide = &ast.BinaryExpression{
			Operator:     token.OperatorValue(operator),
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			NodePosition: parser.createPosition(beginOffset),
		}
	}
}

func (parser *Parser) parseConstructor() (*ast.MethodCall, error) {
	typeName, err := parser.ParseTypeName()
	if err != nil {
		return nil, err
	}
	return parser.parseMethodCallOnNode(typeName)
}

func (parser *Parser) parseCreateExpression() (ast.Node, error) {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.CreateKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	constructor, err := parser.parseConstructor()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	return &ast.CreateExpression{
		NodePosition: parser.createPosition(beginOffset),
		Constructor: constructor,
	}, nil
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parser *Parser) ParseUnaryExpression() (ast.Node, error) {
	beginOffset := parser.offset()
	operatorToken := parser.token()
	if !token.IsOperatorOrOperatorKeywordToken(operatorToken) {
		return parser.ParseOperation()
	}
	if token.KeywordValue(operatorToken) == token.CreateKeyword {
		return parser.parseCreateExpression()
	}
	operator := token.OperatorValue(operatorToken)
	if !operator.IsUnaryOperator() {
		return parser.ParseOperation()
	}
	parser.advance()
	operand, err := parser.ParseUnaryExpression()
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpression{
		Operator: operator,
		Operand:  operand,
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}

// ParseMethodCall parses the call to a method.
func (parser *Parser) parseMethodCallOnNode(method ast.Node) (*ast.MethodCall, error) {
	beginOffset := parser.offset()
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
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}

// parseArgumentList parses the arguments of a MethodCall.
func (parser *Parser) parseArgumentList() ([]ast.Node, error) {
	if token.OperatorValue(parser.token()) == token.RightParenOperator {
		parser.advance()
		return []ast.Node{}, nil
	}
	var arguments []ast.Node
	for {
		next, err := parser.ParseExpression()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, next)
		current := parser.token()
		switch token.OperatorValue(current) {
		case token.RightParenOperator:
			parser.advance()
			return arguments, nil
		case token.CommaOperator:
			parser.advance()
			continue
		}
		return arguments, &UnexpectedTokenError{
			Token:    current,
			Expected: "end of method",
		}
	}
}
