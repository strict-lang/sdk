// The expression file contains methods that are parsing expressions. Every method expects
// the first token that it requires to be the current one (parsing.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the parsing.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package parsing

import (
	"fmt"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

func (parsing *Parsing) parseExpression() (syntaxtree2.Node, error) {
	return parsing.parseBinaryExpression(token2.LowPrecedence + 1)
}

func (parsing *Parsing) parseOperand() (syntaxtree2.Node, error) {
	switch last := parsing.token(); {
	case token2.IsIdentifierToken(last):
		return parsing.parseIdentifier()
	case token2.IsStringLiteralToken(last):
		return parsing.parseStringLiteral()
	case token2.IsNumberLiteralToken(last):
		return parsing.parseNumberLiteral()
	case token2.OperatorValue(last) == token2.LeftParenOperator:
		return parsing.completeLeftParenExpression()
	}
	return nil, fmt.Errorf("could not parse operand: %s", parsing.token())
}

func (parsing *Parsing) parseIdentifier() (*syntaxtree2.Identifier, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree2.Identifier{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseStringLiteral() (*syntaxtree2.StringLiteral, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree2.StringLiteral{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseNumberLiteral() (*syntaxtree2.NumberLiteral, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree2.NumberLiteral{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) completeLeftParenExpression() (syntaxtree2.Node, error) {
	parsing.advance()
	parsing.expressionDepth++
	expression, err := parsing.parseExpression()
	if err != nil {
		return expression, err
	}
	parsing.expressionDepth--
	if token2.OperatorValue(parsing.token()) != token2.RightParenOperator {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: token2.RightParenOperator.String(),
		}
	}
	parsing.advance()
	return expression, nil
}

// ParseOperation parses the initial operand and continues to parsing operands on
// that operand, forming a node for another expression.
func (parsing *Parsing) parseOperation() (syntaxtree2.Node, error) {
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
func (parsing *Parsing) parseOperationOnOperand(operand syntaxtree2.Node) (done bool, node syntaxtree2.Node, err error) {
	switch next := parsing.token(); {
	case token2.HasOperatorValue(next, token2.LeftBracketOperator):
		node, err = parsing.parseListSelectExpression(operand)
		return false, node, err
	case token2.HasOperatorValue(next, token2.LeftParenOperator):
		node, err = parsing.parseCallOnNode(operand)
		return false, node, err
	case token2.HasOperatorValue(next, token2.DotOperator):
		node, err = parsing.parseSelectExpression(operand)
		return false, node, err
	}
	return true, operand, nil
}

func (parsing *Parsing) parseListSelectExpression(target syntaxtree2.Node) (syntaxtree2.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token2.LeftBracketOperator); err != nil {
		return nil, err
	}
	index, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	defer parsing.advance()
	if !token2.HasOperatorValue(parsing.token(), token2.RightBracketOperator) {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "] / end of list access",
		}
	}
	return &syntaxtree2.ListSelectExpression{
		Index:        target,
		Target:       index,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseSelectExpression(target syntaxtree2.Node) (syntaxtree2.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token2.DotOperator); err != nil {
		return nil, err
	}
	field, err := parsing.parseOperand()
	if err != nil {
		return nil, err
	}
	return &syntaxtree2.SelectExpression{
		Target:       target,
		Selection:    field,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parsing *Parsing) parseBinaryExpression(requiredPrecedence token2.Precedence) (syntaxtree2.Node, error) {
	beginOffset := parsing.offset()
	leftHandSide, err := parsing.parseUnaryExpression()
	if err != nil {
		return nil, err
	}
	for {
		operator := parsing.token()
		precedence := token2.PrecedenceOfAny(operator)
		if precedence < requiredPrecedence {
			return leftHandSide, nil
		}
		parsing.advance()
		rightHandSide, err := parsing.parseBinaryExpression(precedence)
		if err != nil {
			return leftHandSide, err
		}
		leftHandSide = &syntaxtree2.BinaryExpression{
			Operator:     token2.OperatorValue(operator),
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
}

func (parsing *Parsing) parseConstructor() (*syntaxtree2.CallExpression, syntaxtree2.TypeName, error) {
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, nil, err
	}
	methodCall, err := parsing.parseCallOnNode(typeName)
	return methodCall, typeName, err
}

func (parsing *Parsing) parseCreateExpression() (syntaxtree2.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.CreateKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	constructor, typeName, err := parsing.parseConstructor()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	return &syntaxtree2.CreateExpression{
		NodePosition: parsing.createPosition(beginOffset),
		Constructor:  constructor,
		Type:         typeName,
	}, nil
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parsing *Parsing) parseUnaryExpression() (syntaxtree2.Node, error) {
	beginOffset := parsing.offset()
	operatorToken := parsing.token()
	if token2.KeywordValue(operatorToken) == token2.CreateKeyword {
		return parsing.parseCreateExpression()
	}
	if !token2.IsOperatorOrOperatorKeywordToken(operatorToken) {
		return parsing.parseOperation()
	}
	operator := token2.OperatorValue(operatorToken)
	if !operator.IsUnaryOperator() {
		return parsing.parseOperation()
	}
	parsing.advance()
	operand, err := parsing.parseUnaryExpression()
	if err != nil {
		return nil, err
	}
	return &syntaxtree2.UnaryExpression{
		Operator:     operator,
		Operand:      operand,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseMethodCall parses the call to a method.
func (parsing *Parsing) parseCallOnNode(method syntaxtree2.Node) (*syntaxtree2.CallExpression, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token2.LeftParenOperator); err != nil {
		return &syntaxtree2.CallExpression{}, err
	}
	arguments, err := parsing.parseArgumentList()
	if err != nil {
		return &syntaxtree2.CallExpression{}, err
	}
	return &syntaxtree2.CallExpression{
		Arguments:    arguments,
		Method:       method,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseCallArgument() (*syntaxtree2.CallArgument, error) {
	beginOffset := parsing.offset()
	var argument syntaxtree2.CallArgument
	if token2.IsIdentifierToken(parsing.token()) &&
		token2.HasOperatorValue(parsing.peek(), token2.AssignOperator) {
		fmt.Println("Label: ", parsing.token().Value())
		argument.Label = parsing.token().Value()
		parsing.advance()
		parsing.advance()
		fmt.Println("Looking at: ", parsing.token().Value())
	}
	value, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	argument.Value = value
	argument.NodePosition = parsing.createPosition(beginOffset)
	return &argument, nil
}

// parseArgumentList parses the arguments of a CallExpression.
func (parsing *Parsing) parseArgumentList() ([]*syntaxtree2.CallArgument, error) {
	if token2.OperatorValue(parsing.token()) == token2.RightParenOperator {
		parsing.advance()
		return []*syntaxtree2.CallArgument{}, nil
	}
	var arguments []*syntaxtree2.CallArgument
	for {
		argument, err := parsing.parseCallArgument()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, argument)
		current := parsing.token()

		switch token2.OperatorValue(current) {
		case token2.RightParenOperator:
			parsing.advance()
			return arguments, nil
		case token2.CommaOperator:
			parsing.advance()
			continue
		}
		return arguments, &UnexpectedTokenError{
			Token:    current,
			Expected: "end of method call",
		}
	}
}
