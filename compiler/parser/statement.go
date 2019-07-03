package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

func (parser *Parser) reportInvalidStatement() {
	offset := parser.tokens.Peek().Position().Begin
	lineIndex := parser.linemap.LineAtOffset(offset)
	parser.reportError(&InvalidStatementError {
		LineIndex: lineIndex,
	})
}

func (parser *Parser) checkLookingAtEndOfStatement() {
	if next := parser.tokens.Pull(); !token.IsEndOfStatementToken(next) {
		parser.reportError(&UnexpectedTokenError{
			Token: next,
			Expected: "end of statement",
		})
	}
}

// ParseIfStatement parses a conditional statement and it's optional else-clause.
func (parser *Parser) ParseIfStatement() ast.Node {
	if err := parser.expectKeyword(token.IfKeyword); err != nil {
		parser.reportError(err)
		return &ast.InvalidStatement{}
	}
	parser.tokens.Pull()
	condition, err := parser.ParseExpression()
	if err != nil {
		parser.reportError(err)
		return &ast.InvalidStatement{}
	}
	parser.checkLookingAtEndOfStatement()
	body := parser.ParseStatementBlock()

	if !token.HasKeywordValue(parser.tokens.Pull(), token.ElseKeyword) {
		return &ast.ConditionalStatement{
			Condition: condition,
			Body: body,
		}
	}
	parser.checkLookingAtEndOfStatement()
	elseBody := parser.ParseStatementBlock()
	return &ast.ConditionalStatement{
		Condition: condition,
		Body: body,
		Else: elseBody,
	}
}

func (parser *Parser) ParseForStatement() ast.Node {
	if err := parser.expectKeyword(token.ForKeyword); err != nil {
		parser.reportError(err)
		return &ast.InvalidStatement{}
	}
	initializerBeginToken := parser.tokens.Pull()
	if token.IsIdentifierToken(initializerBeginToken) {
		if token.HasKeywordValue(parser.tokens.Peek(), token.FromKeyword) {
			return parser.completeFromToStatement()
		}
	}
	return parser.completeForEachStatement()
}

func (parser *Parser) completeForEachStatement() ast.Node {
	field := parser.tokens.Last()
	if !token.IsIdentifierToken(field) {
		parser.reportError(&UnexpectedTokenError{
			Token: field,
			Expected: "identifier",
		})
		return &ast.InvalidStatement{}
	}
	if err :=	parser.skipKeyword(token.FromKeyword); err != nil {
		parser.reportError(err)
		return &ast.InvalidStatement{}
	}
	from, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if err := parser.skipKeyword(token.ToKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	to, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	body := parser.ParseStatementBlock()
	return &ast.FromToLoopStatement{
		Field: ast.NewIdentifier(field.Value()),
		From: from,
		To: to,
		Body: body,
	}
}

func (parser *Parser) completeFromToStatement() ast.Node {

}

// ParseYieldStatement parses a 'yield' statement. Yield statements add an
// element to an implicitly created list, which is returned by the method.
// Any kind of expression can be yielded.
func (parser *Parser) ParseYieldStatement() ast.Node {
	if err := parser.expectKeyword(token.YieldKeyword); err != nil {
		parser.reportError(err)
		return &ast.InvalidStatement{}
	}
	parser.tokens.Pull()
	rightHandSide, err := parser.ParseRightHandSide()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return &ast.YieldStatement{
		Value: rightHandSide,
	}
}

// ParseReturnStatement parses a 'return' statement. Return statements are
// part of the control flow and used to either return a value from a method
// or to end the method call within a branch, resulting in the remaining
// instructions to be ignored. The ReturnStatement is always the last statement
// within a StatementSequence / Branch.
func (parser *Parser) ParseReturnStatement() ast.Node {
	if err := parser.expectKeyword(token.ReturnKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}

	nextToken := parser.tokens.Pull()
	if token.IsEndOfStatementToken(nextToken) {
		return &ast.ReturnStatement{}
	}
	rightHandSide, err := parser.ParseRightHandSide()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return &ast.ReturnStatement{
		Value: rightHandSide,
	}
}

// keywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parser *Parser) keywordStatementParser(keyword token.Keyword) (func() ast.Node, bool) {
	switch keyword {
	case token.IfKeyword:
		return parser.ParseIfStatement, true
	case token.ForKeyword:
		return parser.ParseForStatement, true
	case token.YieldKeyword:
		return parser.ParseYieldStatement, true
	case token.ReturnKeyword:
		return parser.ParseReturnStatement, true
	}
	return nil, false
}

// ParseKeywordStatement parses a statement that starts with a keyword.
func (parser *Parser) ParseKeywordStatement(keyword token.Keyword) ast.Node {
	function, ok := parser.keywordStatementParser(keyword)
	if ok {
		return function()
	}
	parser.reportError(&UnexpectedTokenError{
		Token:    parser.tokens.Peek(),
		Expected: "statement begin",
	})
	return &ast.InvalidStatement{}
}

// parseAssignStatement completes the parsing of a instruction and produces an
// AssignStatement node. Assignments also include those using the Add,Sub,Mod,...Assign
// operators. This method requires that a leftHandSide expression has already been parsed.
func (parser *Parser) parseAssignStatement(operator token.Operator, leftHandSide ast.Node) (ast.Node, error) {
	parser.tokens.Pull()
	rightHandSide, err := parser.ParseRightHandSide()
	if err != nil {
		return &ast.InvalidStatement{}, err
	}
	return &ast.AssignStatement{
		Target: leftHandSide,
		Value: rightHandSide,
		Operator: operator,
	}, nil
}

// ParseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parser *Parser) ParseInstructionStatement() (ast.Node, error) {
	leftHandSide, err := parser.ParseLeftHandSide()
	if err != nil {
		return nil, err
	}
	nextToken := parser.tokens.Pull()
	switch operator := token.OperatorValue(nextToken); {
	case operator.IsAssign():
		return parser.parseAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
		return &ast.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token.DecrementOperator:
		return &ast.DecrementStatement{Operand: leftHandSide}, nil
	}
	return &ast.InvalidStatement{}, &UnexpectedTokenError{
		Token: nextToken,
		Expected: "operator",
	}
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parser *Parser) ParseStatement() ast.Node {
	switch peek := parser.tokens.Peek(); {
	case token.IsKeywordToken(peek):
		return parser.ParseKeywordStatement(peek.(*token.KeywordToken).Keyword)
	case token.IsIdentifierToken(peek),
		token.IsOperatorToken(peek),
		token.IsLiteralToken(peek):

		statement, err := parser.ParseInstructionStatement()
		if err != nil {
			return parser.createInvalidStatement(err)
		}
		nextToken := parser.tokens.Pull()
		if !token.IsEndOfStatementToken(nextToken) {
			parser.reportError(&UnexpectedTokenError{
				Token: nextToken,
				Expected: "end of statement",
			})
		}
		return statement
	default:
		parser.reportInvalidStatement()
		return &ast.InvalidStatement{}
	}
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parser *Parser) ParseStatementSequence() []ast.Node {
	var statements []ast.Node
	for {
		next := parser.tokens.Pull()
		if next.Indent() != parser.block.Indent {
			break
		}
		statements = append(statements, parser.ParseStatement())
	}
	return statements
}

func (parser *Parser) ParseStatementBlock() ast.Node {
	return nil
}
