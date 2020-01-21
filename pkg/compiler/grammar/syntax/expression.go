// The expression file contains methods that are grammar expressions. Every method expects
// the first token that it requires to be the current one (grammar.token()) it responsible
// to advance all tokens so that the next method can directly continue without having to
// call the grammar.advance() method itself. This is done because developers should always
// know what the current token is, to prevent bugs.

package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

func (parsing *Parsing) parseExpression() tree.Expression {
	return parsing.parseAnyExpression()
}

func (parsing *Parsing) parseConditionalExpression() tree.Expression {
	return parsing.parseBinaryExpression(token.InitialConditionalPrecedence)
}

func (parsing *Parsing) parseAnyExpression() tree.Expression {
	return parsing.parseBinaryExpression(token.InitialPrecedence)
}

// ParseUnaryExpression parses a unary expression. Unary expressions are
// operations with only one operand (arity of one). An example of a unary
// expression is the negation '!(expression)'. The single operand may be
// any kind of expression, including another unary expression.
func (parsing *Parsing) parseUnaryExpression() tree.Expression {
	parsing.beginStructure(tree.UnaryExpressionNodeKind)
	// TODO: Write tests to find out if other methods change the kind
	defer parsing.completeStructure(tree.UnaryExpressionNodeKind)
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
		Region:   parsing.createRegionOfCurrentStructure(),
	}
}

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

func (parsing *Parsing) parseOperand() tree.Expression {
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
	parsing.throwError(&diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "operand",
			Received: parsing.token().Value(),
		},
		CommonReasons: []string{
			"The expression is not a valid operand",
			"An unsupported operation or operator is used",
		},
	})
}

func (parsing *Parsing) parseIdentifier() *tree.Identifier {
	defer parsing.advance()
	return parsing.expectAnyIdentifier()
}

func (parsing *Parsing) parseStringLiteral() *tree.StringLiteral {
	parsing.beginStructure(tree.StringLiteralNodeKind)
	literalToken := parsing.pullToken()
	return &tree.StringLiteral{
		Value:  literalToken.Value(),
		Region: parsing.completeStructure(tree.StringLiteralNodeKind),
	}
}

func (parsing *Parsing) parseNumberLiteral() *tree.NumberLiteral {
	parsing.beginStructure(tree.IdentifierNodeKind)
	literalToken := parsing.pullToken()
	return &tree.NumberLiteral{
		Value:  literalToken.Value(),
		Region: parsing.completeStructure(tree.IdentifierNodeKind),
	}
}

func (parsing *Parsing) completeLeftParenExpression() tree.Expression {
	parsing.advance()
	parsing.expressionDepth++
	expression := parsing.parseExpression()
	parsing.expressionDepth--
	parsing.expectEndOfLeftParenExpression()
	return expression
}

func (parsing *Parsing) expectEndOfLeftParenExpression() {
	if token.OperatorValue(parsing.pullToken()) != token.RightParenOperator {
		parsing.throwError(&diagnostic.RichError{
			Error: &diagnostic.UnexpectedTokenError{
				Expected: ")",
				Received: parsing.token().Value(),
			},
			CommonReasons: []string{
				"An right paren is missing",
				"The expression is invalid",
			},
		})
	}
}

// ParseOperation parses the initial operand and continues to grammar operands on
// that operand, forming a node for another expression.
func (parsing *Parsing) parseOperation() tree.Expression {
	operand := parsing.parseOperand()
	return parsing.parseOperationsOnOperand(operand)
}

func (parsing *Parsing) parseOperationOrAssign(node tree.Node) tree.Node {
	if token.IsOperatorToken(parsing.token()) {
		operator := token.OperatorValue(parsing.token())
		parsing.advance()
		return parsing.completeAssignStatement(operator, node)
	}
	return node
}

func (parsing *Parsing) parseOperationsOnOperand(operand tree.Expression) tree.Expression {
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
func (parsing *Parsing) parseOperationOnOperand(
	operand tree.Expression) (node tree.Expression, done bool) {

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

func (parsing *Parsing) parseListSelectExpression(
	target tree.Expression) *tree.ListSelectExpression {

	parsing.beginStructure(tree.ListSelectExpressionNodeKind)
	parsing.skipOperator(token.LeftBracketOperator)
	index := parsing.parseExpression()
	parsing.expectEndOfListSelect()
	return &tree.ListSelectExpression{
		Index:  index,
		Target: target,
		Region: parsing.completeStructure(tree.ListSelectExpressionNodeKind),
	}
}

func (parsing *Parsing) expectEndOfListSelect() {
	if !token.HasOperatorValue(parsing.pullToken(), token.RightBracketOperator) {
		parsing.throwError(newEndOfListSelectError(parsing.token()))
	}
}

func newEndOfListSelectError(received token.Token) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "end of list selection",
			Received: received.Value(),
		},
		CommonReasons: []string{
			"A ListSelectExpression is not closed in the line it was opened",
			"The index expression is invalid",
		},
	}
}

func (parsing *Parsing) parseFieldSelectExpression(target tree.Expression) *tree.FieldSelectExpression {
	parsing.beginStructure(tree.FieldSelectExpressionNodeKind)
	parsing.skipOperator(token.DotOperator)
	field := parsing.parseOperand()
	return &tree.FieldSelectExpression{
		Target:    target,
		Selection: field,
		Region:    parsing.completeStructure(tree.FieldSelectExpressionNodeKind),
	}
}

// TODO: Get rid of special syntax for constructor calls.
func (parsing *Parsing) parseConstructorCall() (*tree.CallExpression, tree.TypeName) {
	typeName := parsing.parseTypeName()
	methodCall := parsing.parseCallOnNode(&tree.Identifier{
		Value: typeName.BaseName(),
	})
	return methodCall, typeName
}

func (parsing *Parsing) parseCreateExpression() tree.Expression{
	parsing.beginStructure(tree.CreateExpressionNodeKind)
	parsing.skipKeyword(token.CreateKeyword)
	constructor, typeName := parsing.parseConstructorCall()
	return &tree.CreateExpression{
		Call:   constructor,
		Type:   typeName,
		Region: parsing.completeStructure(tree.CreateExpressionNodeKind),
	}
}

// ParseMethodCall parses the call to a method.
func (parsing *Parsing) parseCallOnNode(method tree.Expression) *tree.CallExpression {
	parsing.beginStructure(tree.CallExpressionNodeKind)
	parsing.skipOperator(token.LeftParenOperator)
	arguments := parsing.parseArgumentList()
	return &tree.CallExpression{
		Arguments: arguments,
		Target:    method,
		Region:    parsing.completeStructure(tree.CallExpressionNodeKind),
	}
}

func (parsing *Parsing) parseCallArgument() *tree.CallArgument {
	parsing.beginStructure(tree.CallArgumentNodeKind)
	var argument tree.CallArgument
	if token.IsIdentifierToken(parsing.token()) &&
		token.HasOperatorValue(parsing.peek(), token.AssignOperator) {
		argument.Label = parsing.token().Value()
		parsing.advance()
		parsing.advance()
	}
	value := parsing.parseExpression()
	argument.Value = value
	argument.Region = parsing.completeStructure(tree.CallArgumentNodeKind)
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
	operator := token.OperatorValue(parsing.token())
	if operator == token.CommaOperator {
		parsing.advance()
		return
	}
	if operator != token.RightParenOperator {
		parsing.throwExpectedEndOfMethodCallError()
	}
}

func (parsing *Parsing) throwExpectedEndOfMethodCallError() {
	parsing.throwError(&diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "end of call",
			Received: parsing.token().Value(),
		},
		CommonReasons: []string{
			"A method call is not closed in the right line",
			"The argument ist is not properly separated",
			"The last argument expression is invalid",
		},
	})
}

func (parsing *Parsing) isAtEndOfArgumentList() bool {
	return token.HasOperatorValue(parsing.token(), token.RightParenOperator)
}
