// The expression file contains methods that are parsing expressions. Every method expects
// the first token that it requires to be the current one (parsing.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the parsing.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package parsing

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

func (parsing *Parsing) parseExpression() (ast.Node, error) {
	return parsing.parseBinaryExpression(token.LowPrecedence + 1)
}

func (parsing *Parsing) parseOperand() (ast.Node, error) {
	switch last := parsing.token(); {
	case token.IsIdentifierToken(last):
		return parsing.parseIdentifier()
	case token.IsStringLiteralToken(last):
		return parsing.parseStringLiteral()
	case token.IsNumberLiteralToken(last):
		return parsing.parseNumberLiteral()
	case token.OperatorValue(last) == token.LeftParenOperator:
		return parsing.completeLeftParenExpression()
	}
	return nil, fmt.Errorf("could not parse operand: %s", parsing.token())
}

func (parsing *Parsing) parseIdentifier() (*ast.Identifier, error) {
	defer parsing.advance()
	return &ast.Identifier{
		Value:        parsing.token().Value(),
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseStringLiteral() (*ast.StringLiteral, error) {
	defer parsing.advance()
	return &ast.StringLiteral{
		Value:        parsing.token().Value(),
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseNumberLiteral() (*ast.NumberLiteral, error) {
	defer parsing.advance()
	return &ast.NumberLiteral{
		Value:        parsing.token().Value(),
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) completeLeftParenExpression() (ast.Node, error) {
	parsing.advance()
	parsing.expressionDepth++
	expression, err := parsing.parseExpression()
	if err != nil {
		return expression, err
	}
	parsing.expressionDepth--
	if token.OperatorValue(parsing.token()) != token.RightParenOperator {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: token.RightParenOperator.String(),
		}
	}
	parsing.advance()
	return expression, nil
}

// ParseOperation parses the initial operand and continues to parsing operands on
// that operand, forming a node for another expression.
func (parsing *Parsing) parseOperation() (ast.Node, error) {
	operand, err := parsing.parseOperand()
	if err != nil {
		return nil, err
	}
	for {
		done, node, err := parsing.parseOperationOnOperand(operand)
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
func (parsing *Parsing) parseOperationOnOperand(operand ast.Node) (done bool, node ast.Node, err error) {
	switch next := parsing.token(); {
	case token.OperatorValue(next) == token.LeftParenOperator:
		call, err := parsing.parseMethodCallOnNode(operand)
		return false, call, err
	case token.OperatorValue(next) == token.DotOperator:
		selector, err := parsing.parseSelection(operand)
		return false, selector, err
	}
	return true, operand, nil
}

func (parsing *Parsing) parseSelection(operand ast.Node) (ast.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token.DotOperator); err != nil {
		return nil, err
	}
	field, err := parsing.parseOperand()
	if err != nil {
		return nil, err
	}
	return &ast.SelectorExpression{
		Target:       operand,
		Selection:    field,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parsing *Parsing) parseBinaryExpression(requiredPrecedence token.Precedence) (ast.Node, error) {
	beginOffset := parsing.offset()
	leftHandSide, err := parsing.parseUnaryExpression()
	if err != nil {
		return nil, err
	}
	for {
		operator := parsing.token()
		precedence := token.PrecedenceOfAny(operator)
		if precedence < requiredPrecedence {
			return leftHandSide, nil
		}
		parsing.advance()
		rightHandSide, err := parsing.parseBinaryExpression(precedence)
		if err != nil {
			return leftHandSide, err
		}
		leftHandSide = &ast.BinaryExpression{
			Operator:     token.OperatorValue(operator),
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
}

func (parsing *Parsing) parseConstructor() (*ast.MethodCall, ast.TypeName, error) {
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, nil, err
	}
	methodCall, err := parsing.parseMethodCallOnNode(typeName)
	return methodCall, typeName, err
}

func (parsing *Parsing) parseCreateExpression() (ast.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.CreateKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	constructor, typeName, err := parsing.parseConstructor()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	return &ast.CreateExpression{
		NodePosition: parsing.createPosition(beginOffset),
		Constructor:  constructor,
		Type: typeName,
	}, nil
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parsing *Parsing) parseUnaryExpression() (ast.Node, error) {
	beginOffset := parsing.offset()
	operatorToken := parsing.token()
	if token.KeywordValue(operatorToken) == token.CreateKeyword {
		return parsing.parseCreateExpression()
	}
	if !token.IsOperatorOrOperatorKeywordToken(operatorToken) {
		return parsing.parseOperation()
	}
	operator := token.OperatorValue(operatorToken)
	if !operator.IsUnaryOperator() {
		return parsing.parseOperation()
	}
	parsing.advance()
	operand, err := parsing.parseUnaryExpression()
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpression{
		Operator:     operator,
		Operand:      operand,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseMethodCall parses the call to a method.
func (parsing *Parsing) parseMethodCallOnNode(method ast.Node) (*ast.MethodCall, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token.LeftParenOperator); err != nil {
		return &ast.MethodCall{}, err
	}
	arguments, err := parsing.parseArgumentList()
	if err != nil {
		return &ast.MethodCall{}, err
	}
	return &ast.MethodCall{
		Arguments:    arguments,
		Method:       method,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// parseArgumentList parses the arguments of a MethodCall.
func (parsing *Parsing) parseArgumentList() ([]ast.Node, error) {
	if token.OperatorValue(parsing.token()) == token.RightParenOperator {
		parsing.advance()
		return []ast.Node{}, nil
	}
	var arguments []ast.Node
	for {
		next, err := parsing.parseExpression()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, next)
		current := parsing.token()
		switch token.OperatorValue(current) {
		case token.RightParenOperator:
			parsing.advance()
			return arguments, nil
		case token.CommaOperator:
			parsing.advance()
			continue
		}
		return arguments, &UnexpectedTokenError{
			Token:    current,
			Expected: "end of method",
		}
	}
}
