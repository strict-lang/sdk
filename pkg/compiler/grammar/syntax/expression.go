// The expression file contains methods that are grammar expressions. Every method expects
// the first token that it requires to be the current one (grammar.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the grammar.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (parsing *Parsing) parseExpression() tree.Node {
	return parsing.parseBinaryExpression(token.LowPrecedence + 1)
}

func (parsing *Parsing) parseOperand() tree.Node {
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
	parsing.throwInvalidOperandError()
	return nil
}

func (parsing *Parsing) throwInvalidOperandError() {
	err := fmt.Errorf("could not parse operand: %s", parsing.token())
	parsing.throwError(err)
}

func (parsing *Parsing) parseIdentifier() *tree.Identifier {
	identifier := parsing.pullToken()
	return &tree.Identifier{
		Value:  identifier.Value(),
		Region: parsing.createRegionFromToken(identifier),
	}
}

func (parsing *Parsing) parseStringLiteral() *tree.StringLiteral {
	literalToken := parsing.pullToken()
	return &tree.StringLiteral{
		Value:  literalToken.Value(),
		Region: parsing.createRegionFromToken(literalToken),
	}
}

func (parsing *Parsing) parseNumberLiteral() *tree.NumberLiteral {
	literalToken := parsing.pullToken()
	return &tree.NumberLiteral{
		Value:  literalToken.Value(),
		Region: parsing.createRegionFromToken(literalToken),
	}
}

func (parsing *Parsing) completeLeftParenExpression() tree.Node {
	parsing.advance()
	parsing.expressionDepth++
	expression := parsing.parseExpression()
	parsing.expressionDepth--
	parsing.expectEndOfLeftParenExpression()
	return expression
}

func (parsing *Parsing) expectEndOfLeftParenExpression() {
	if token.OperatorValue(parsing.pullToken()) != token.RightParenOperator {
		parsing.throwError(&UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: token.RightParenOperator.String(),
		})
	}
}

// ParseOperation parses the initial operand and continues to grammar operands on
// that operand, forming a node for another expression.
func (parsing *Parsing) parseOperation() tree.Node {
	operand := parsing.parseOperand()
	return parsing.parseOperationsOnOperand(operand)
}

func (parsing *Parsing) parseOperationOrAssign(
	node tree.Node) (tree.Node, error) {

	if token.IsOperatorToken(parsing.token()) {
		operator := token.OperatorValue(parsing.token())
		parsing.advance()
		return parsing.parseAssignStatement(operator, node)
	}
	return node, nil
}

func (parsing *Parsing) parseOperationsOnOperand(operand tree.Node) tree.Node {
	for {
		if node, done := parsing.parseOperationOnOperand(operand); !done {
			operand = node
		} else {
			return operand
		}
	}
}

// ParseOperationOnOperand parses an operation on an operand that has already
// been parsed. It is called by the ParseOperand method.
func (parsing *Parsing) parseOperationOnOperand(operand tree.Node) (node tree.Node, done bool) {
	switch next := parsing.token(); {
	case token.HasOperatorValue(next, token.LeftParenOperator):
		return parsing.parseCallOnNode(operand), false
	case token.HasOperatorValue(next, token.LeftBracketOperator):
		return parsing.parseListSelectExpression(operand), false
	case token.HasOperatorValue(next, token.DotOperator):
		return parsing.parseFieldSelectExpression(operand), false
	}
	return nil, true
}

func (parsing *Parsing) parseListSelectExpression(target tree.Node) *tree.ListSelectExpression {
	beginOffset := parsing.offset()
	parsing.skipOperator(token.LeftBracketOperator)
	index := parsing.parseExpression()
	parsing.expectEndOfListSelect()
	return &tree.ListSelectExpression{
		Index:  index,
		Target: target,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) expectEndOfListSelect() {
	if !token.HasOperatorValue(parsing.pullToken(), token.RightBracketOperator) {
		parsing.throwError(&UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "] / end of list access",
		})
	}
}

func (parsing *Parsing) parseFieldSelectExpression(target tree.Node) *tree.FieldSelectExpression {
	beginOffset := parsing.offset()
	parsing.skipOperator(token.DotOperator)
	field := parsing.parseOperand()
	return &tree.FieldSelectExpression{
		Target:    target,
		Selection: field,
		Region:    parsing.createRegion(beginOffset),
	}
}

// ParseBinaryExpression parses a binary expression. Binary expressions are
// operations with two operands. Strict uses the infix notation, therefor
// binary expressions have a left-hand-side and right-hand-side operand and
// the operator in between. The operands can be any kind of expression.
// Example: 'a + b' or '(1 + 2) + 3'
func (parsing *Parsing) parseBinaryExpression(requiredPrecedence token.Precedence) tree.Node {
	expression, _ := parsing.parseBinaryExpressionRecursive(requiredPrecedence)
	return expression
}

func (parsing *Parsing) parseBinaryExpressionRecursive(
	requiredPrecedence token.Precedence) (tree.Node, bool) {

	beginOffset := parsing.offset()
	leftHandSide := parsing.parseUnaryExpression()
	for {
		operator := parsing.token()
		precedence := token.PrecedenceOfAny(operator)
		if precedence < requiredPrecedence {
			return leftHandSide, false
		}
		parsing.advance()
		rightHandSide, success := parsing.parseBinaryExpressionRecursive(precedence + 1)
		if !success {
			return leftHandSide, true
		}
		leftHandSide = &tree.BinaryExpression{
			Operator:     token.OperatorValue(operator),
			LeftOperand:  leftHandSide,
			RightOperand: rightHandSide,
			Region:       parsing.createRegion(beginOffset),
		}
	}
}

func (parsing *Parsing) parseConstructorCall() (*tree.CallExpression, tree.TypeName) {
	typeName := parsing.parseTypeName()
	methodCall := parsing.parseCallOnNode(typeName)
	return methodCall, typeName
}

func (parsing *Parsing) parseCreateExpression() tree.Node {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.CreateKeyword)
	constructor, typeName := parsing.parseConstructorCall()
	return &tree.CreateExpression{
		Call:   constructor,
		Type:   typeName,
		Region: parsing.createRegion(beginOffset),
	}
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parsing *Parsing) parseUnaryExpression() tree.Node {
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
	operand := parsing.parseUnaryExpression()
	return &tree.UnaryExpression{
		Operator: operator,
		Operand:  operand,
		Region:   parsing.createRegion(beginOffset),
	}
}

// ParseMethodCall parses the call to a method.
func (parsing *Parsing) parseCallOnNode(method tree.Node) *tree.CallExpression {
	beginOffset := parsing.offset()
	parsing.skipOperator(token.LeftParenOperator)
	arguments := parsing.parseArgumentList()
	return &tree.CallExpression{
		Arguments: arguments,
		Target:    method,
		Region:    parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseCallArgument() *tree.CallArgument {
	beginOffset := parsing.offset()
	var argument tree.CallArgument
	if token.IsIdentifierToken(parsing.token()) &&
		token.HasOperatorValue(parsing.peek(), token.AssignOperator) {
		argument.Label = parsing.token().Value()
		parsing.advance()
		parsing.advance()
	}
	value := parsing.parseExpression()
	argument.Value = value
	argument.Region = parsing.createRegion(beginOffset)
	return &argument
}

// parseArgumentList parses the arguments of a CallExpression.
func (parsing *Parsing) parseArgumentList() tree.CallArgumentList {
	if parsing.isAtEndOfArgumentList() {
		parsing.advance()
		return tree.CallArgumentList{}
	}
	return parsing.parseNonEmptyArgumentList()
}

func (parsing *Parsing) parseNonEmptyArgumentList() (arguments tree.CallArgumentList) {
	for !parsing.isAtEndOfArgumentList() {
		arguments = append(arguments, parsing.parseCallArgument())
		parsing.consumeTokenAfterArgument()
	}
	return arguments
}

func (parsing *Parsing) consumeTokenAfterArgument() {
	if !token.HasOperatorValue(parsing.pullToken(), token.CommaOperator) {
		parsing.throwExpectedEndOfMethodCallError()
	}
}

func (parsing *Parsing) throwExpectedEndOfMethodCallError() {
	parsing.throwError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "end of method call",
	})
}

func (parsing *Parsing) isAtEndOfArgumentList() bool {
	return token.HasOperatorValue(parsing.token(), token.RightParenOperator)
}
