// The expression file contains methods that are parsing expressions. Every method expects
// the first token that it requires to be the current one (parsing.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the parsing.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package parsing

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

func (parsing *Parsing) parseExpression() (syntaxtree.Node, error) {
	return parsing.parseBinaryExpression(token.LowPrecedence + 1)
}

func (parsing *Parsing) parseOperand() (syntaxtree.Node, error) {
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

func (parsing *Parsing) parseIdentifier() (*syntaxtree.Identifier, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree.Identifier{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseStringLiteral() (*syntaxtree.StringLiteral, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree.StringLiteral{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) parseNumberLiteral() (*syntaxtree.NumberLiteral, error) {
	value := parsing.token().Value()
	parsing.advance()
	return &syntaxtree.NumberLiteral{
		Value:        value,
		NodePosition: parsing.createTokenPosition(),
	}, nil
}

func (parsing *Parsing) completeLeftParenExpression() (syntaxtree.Node, error) {
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
func (parsing *Parsing) parseOperation() (syntaxtree.Node, error) {
	operand, err := parsing.parseOperand()
	if err != nil {
		return nil, err
	}
	return parsing.parseOperationsOnOperand(operand)
}

func (parsing *Parsing) parseOperationOrAssign(
	node syntaxtree.Node) (syntaxtree.Node, error) {

	if token.IsOperatorToken(parsing.token()) {
		operator := token.OperatorValue(parsing.token())
		parsing.advance()
		return parsing.parseAssignStatement(operator, node)
	}
	return node, nil
}

func (parsing *Parsing) parseOperationsOnOperand(operand syntaxtree.Node) (syntaxtree.Node, error) {
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
func (parsing *Parsing) parseOperationOnOperand(operand syntaxtree.Node) (done bool, node syntaxtree.Node, err error) {
	switch next := parsing.token(); {
	case token.HasOperatorValue(next, token.LeftBracketOperator):
		node, err = parsing.parseListSelectExpression(operand)
		return false, node, err
	case token.HasOperatorValue(next, token.LeftParenOperator):
		node, err = parsing.parseCallOnNode(operand)
		return false, node, err
	case token.HasOperatorValue(next, token.DotOperator):
		node, err = parsing.parseSelectExpression(operand)
		return false, node, err
	}
	return true, operand, nil
}

func (parsing *Parsing) parseListSelectExpression(target syntaxtree.Node) (syntaxtree.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token.LeftBracketOperator); err != nil {
		return nil, err
	}
	index, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	defer parsing.advance()
	if !token.HasOperatorValue(parsing.token(), token.RightBracketOperator) {
		return nil, &UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "] / end of list access",
		}
	}
	return &syntaxtree.ListSelectExpression{
		Index:        index,
		Target:       target,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseSelectExpression(target syntaxtree.Node) (syntaxtree.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token.DotOperator); err != nil {
		return nil, err
	}
	field, err := parsing.parseOperand()
	if err != nil {
		return nil, err
	}
	return &syntaxtree.SelectExpression{
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
func (parsing *Parsing) parseBinaryExpression(requiredPrecedence token.Precedence) (syntaxtree.Node, error) {
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
		rightHandSide, err := parsing.parseBinaryExpression(precedence + 1)
		if err != nil {
			return leftHandSide, err
		}
		leftHandSide = &syntaxtree.BinaryExpression{
			Operator:     token.OperatorValue(operator),
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
}

func (parsing *Parsing) parseConstructor() (*syntaxtree.CallExpression, syntaxtree.TypeName, error) {
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, nil, err
	}
	methodCall, err := parsing.parseCallOnNode(typeName)
	return methodCall, typeName, err
}

func (parsing *Parsing) parseCreateExpression() (syntaxtree.Node, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.CreateKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	constructor, typeName, err := parsing.parseConstructor()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	return &syntaxtree.CreateExpression{
		NodePosition: parsing.createPosition(beginOffset),
		Constructor:  constructor,
		Type:         typeName,
	}, nil
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parsing *Parsing) parseUnaryExpression() (syntaxtree.Node, error) {
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
	return &syntaxtree.UnaryExpression{
		Operator:     operator,
		Operand:      operand,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseMethodCall parses the call to a method.
func (parsing *Parsing) parseCallOnNode(method syntaxtree.Node) (*syntaxtree.CallExpression, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipOperator(token.LeftParenOperator); err != nil {
		return &syntaxtree.CallExpression{}, err
	}
	arguments, err := parsing.parseArgumentList()
	if err != nil {
		return &syntaxtree.CallExpression{}, err
	}
	return &syntaxtree.CallExpression{
		Arguments:    arguments,
		Method:       method,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseCallArgument() (*syntaxtree.CallArgument, error) {
	beginOffset := parsing.offset()
	var argument syntaxtree.CallArgument
	if token.IsIdentifierToken(parsing.token()) &&
		token.HasOperatorValue(parsing.peek(), token.AssignOperator) {
		argument.Label = parsing.token().Value()
		parsing.advance()
		parsing.advance()
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
func (parsing *Parsing) parseArgumentList() ([]*syntaxtree.CallArgument, error) {
	if token.OperatorValue(parsing.token()) == token.RightParenOperator {
		parsing.advance()
		return []*syntaxtree.CallArgument{}, nil
	}
	var arguments []*syntaxtree.CallArgument
	for {
		argument, err := parsing.parseCallArgument()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, argument)
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
			Expected: "end of method call",
		}
	}
}
