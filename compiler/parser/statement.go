package parser

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// ParseIfStatement parses a conditional statement and it's optional else-clause.
func (parser *Parser) ParseIfStatement() ast.Node {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.IfKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	condition, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if !token.HasKeywordValue(parser.token(), token.ElseKeyword) {
		return &ast.ConditionalStatement{
			Condition: condition,
			Consequence:      body,
			NodePosition: parser.createPosition(beginOffset),
		}
	}
	parser.advance()
	elseBody, err := parser.parseElseIfOrBlock()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	return &ast.ConditionalStatement{
		Condition: condition,
		Consequence:      body,
		Alternative:      elseBody,
		NodePosition: parser.createPosition(beginOffset),
	}
}

func (parser *Parser) parseElseIfOrBlock() (ast.Node, error) {
	if token.HasKeywordValue(parser.token(), token.IfKeyword) {
		return parser.ParseIfStatement(), nil
	}
	parser.skipEndOfStatement()
	return parser.ParseStatementBlock()
}

// ParseForStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parser *Parser) ParseForStatement() ast.Node {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.ForKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	initializerBeginToken := parser.token()
	if token.IsIdentifierToken(initializerBeginToken) {
		if token.HasKeywordValue(parser.peek(), token.FromKeyword) {
			return parser.completeFromToStatement(beginOffset)
		}
	}
	return parser.completeForEachStatement(beginOffset)
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops parsing.
func (parser *Parser) completeForEachStatement(beginOffset source.Offset) ast.Node {
	field, err := parser.expectAnyIdentifier()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.advance()
	if err := parser.skipKeyword(token.InKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	value, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	return &ast.ForEachLoopStatement{
		Field:  field,
		Enumeration: value,
		Body:   body,
		NodePosition: parser.createPosition(beginOffset),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops parsing.
func (parser *Parser) completeFromToStatement(beginOffset source.Offset) ast.Node {
	field, err := parser.expectAnyIdentifier()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.advance()
	if err := parser.skipKeyword(token.FromKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	from, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if err := parser.skipKeyword(token.ToKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	to, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	return &ast.RangedLoopStatement{
		ValueField: field,
		InitialValue:  from,
		EndValue:    to,
		Body:  body,
		NodePosition: parser.createPosition(beginOffset),
	}
}

// ParseYieldStatement parses a 'yield' statement. Yield statements add an
// element to an implicitly created list, which is returned by the method.
// Any kind of expression can be yielded.
func (parser *Parser) ParseYieldStatement() ast.Node {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.YieldKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.skipEndOfStatement()
	return &ast.YieldStatement{
		Value: rightHandSide,
		NodePosition: parser.createPosition(beginOffset),
	}
}

// ParseReturnStatement parses a 'return' statement. Return statements are
// part of the control flow and used to either return a value from a method
// or to end the method call within a branch, resulting in the remaining
// instructions to be ignored. The ReturnStatement is always the last statement
// within a StatementSequence / Branch.
func (parser *Parser) ParseReturnStatement() ast.Node {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.ReturnKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	if token.IsEndOfStatementToken(parser.token()) {
		parser.advance()
		return &ast.ReturnStatement{
			NodePosition: parser.createPosition(beginOffset),
		}
	}
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	parser.skipEndOfStatement()
	return &ast.ReturnStatement{
		Value: rightHandSide,
		NodePosition: parser.createPosition(beginOffset),
	}
}

func (parser *Parser) parseNestedMethodDeclaration() ast.Node {
	beginPosition := parser.offset()
	method, err := parser.ParseMethodDeclaration()
	if err != nil {
		return parser.createInvalidStatement(beginPosition, err)
	}
	return method
}

func (parser *Parser) ParseImportStatement() ast.Node {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.ImportKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	path := parser.token()
	if !token.IsStringLiteralToken(path) {
		return parser.createInvalidStatement(beginOffset, &UnexpectedTokenError{
			Expected: "Path",
			Token: path,
		})
	}
	parser.advance()
	if !token.HasKeywordValue(parser.token(), token.AsKeyword) {
		parser.skipEndOfStatement()
		return &ast.ImportStatement{
			Path: path.Value(),
			NodePosition: parser.createPosition(beginOffset),
		}
	}
	parser.advance()
	aliasOffset := parser.offset()
	alias, err := parser.parseImportAlias()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err)
	}
	aliasEnd := parser.offset()
	parser.skipEndOfStatement()
	return &ast.ImportStatement{
		Path: path.Value(),
		Alias: &ast.Identifier{
			Value: alias,
			NodePosition: &offsetPosition{
				begin: aliasOffset,
				end: aliasEnd,
			},
		},
		NodePosition: parser.createPosition(beginOffset),
	}
}

func (parser *Parser) parseImportAlias() (string, error) {
	alias := parser.token()
	if !token.IsIdentifierToken(alias) {
		return "", &UnexpectedTokenError{
			Expected: "Identifier",
			Token: alias,
		}
	}
	parser.advance()
	return alias.Value(), nil
}

var errNoAssign = errors.New("no assign")

func (parser *Parser) parseOptionalAssignValue() (ast.Node, error) {
	if !token.HasOperatorValue(parser.token(), token.AssignOperator) {
		return nil, errNoAssign
	}
	parser.advance()
	return parser.ParseExpression()
}

// keywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parser *Parser) keywordStatementParser(keyword token.Keyword) (func() ast.Node, bool) {
	switch keyword {
	case token.IfKeyword:     return parser.ParseIfStatement, true
	case token.ForKeyword:    return parser.ParseForStatement, true
	case token.YieldKeyword:  return parser.ParseYieldStatement, true
	case token.ReturnKeyword: return parser.ParseReturnStatement, true
	case token.ImportKeyword: return parser.ParseImportStatement, true
	case token.TestKeyword:   return parser.ParseTestStatement, true
	case token.MethodKeyword: return parser.parseNestedMethodDeclaration, true
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
		Token:    parser.token(),
		Expected: "statement begin",
	})
	return &ast.InvalidStatement{}
}

func (parser *Parser) parseAssignStatement(operator token.Operator, leftHandSide ast.Node) (ast.Node, error) {
	beginOffset := parser.offset()
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	parser.skipEndOfStatement()
	return &ast.AssignStatement{
		Target:   leftHandSide,
		Value:    rightHandSide,
		Operator: operator,
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}

func (parser *Parser) ParseTestStatement() (ast.Node, error) {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.TestKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	parser.skipEndOfStatement()
	statements, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	return &ast.TestStatement{
		NodePosition: parser.createPosition(beginOffset),
		MethodName: parser.currentMethodName,
		Statements: statements,
	}, nil
}

func (parser *Parser) ParseAssertStatement() (ast.Node, error) {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.AssertKeyword); err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	expression, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(beginOffset, err), err
	}
	return &ast.AssertStatement{
		NodePosition: parser.createPosition(beginOffset),
		Expression: expression,
	}, nil
}

// ParseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parser *Parser) ParseInstructionStatement() (ast.Node, error) {
	leftHandSide, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}
	switch operator := token.OperatorValue(parser.token()); {
	case operator.IsAssign():
		parser.skipEndOfStatement()
		return parser.parseAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
		parser.advance()
		parser.skipEndOfStatement()
		return &ast.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token.DecrementOperator:
		parser.advance()
		parser.skipEndOfStatement()
		return &ast.DecrementStatement{Operand: leftHandSide}, nil
	}
	parser.advance()
	return &ast.ExpressionStatement{
		Expression: leftHandSide,
	}, nil
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parser *Parser) ParseStatement() ast.Node {
	beginOffset := parser.offset()
	switch current := parser.token(); {
	case token.IsKeywordToken(current):
		return parser.ParseKeywordStatement(token.KeywordValue(current))
	case token.IsIdentifierToken(current):
		fallthrough
	case token.IsOperatorToken(current):
		fallthrough
	case token.IsLiteralToken(current):
		statement, err := parser.ParseInstructionStatement()
		if err != nil {
			return parser.createInvalidStatement(beginOffset, err)
		}
		return statement
	default:
		return parser.createInvalidStatement(beginOffset, ErrInvalidExpression)
	}
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parser *Parser) ParseStatementSequence() []ast.Node {
	var statements []ast.Node
	for {
		expectedIndent := parser.block.Indent
		current := parser.token()
		if token.IsEndOfFileToken(current) {
			break
		}
		if current.Indent() > expectedIndent {
			beginOffset := parser.offset() - 1
			invalid := parser.createInvalidStatement(beginOffset, &InvalidIndentationError{
				Token:    current,
				Expected: fmt.Sprintf("indent level of %d", expectedIndent),
			})
			statements = append(statements, invalid)
			break
		}
		if current.Indent() < expectedIndent {
			break
		}
		statement := parser.ParseStatement()
		if _, ok := statement.(*ast.InvalidStatement); ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

// ParseStatementBlock parses a block of statements.
func (parser *Parser) ParseStatementBlock() (*ast.BlockStatement, error) {
	beginOffset := parser.offset()
	indent := parser.token().Indent()
	if indent < parser.block.Indent {
		return nil, &InvalidIndentationError{
			Token:    parser.token(),
			Expected: "indent bigger than 0",
		}
	}
	parser.openBlock(indent)
	statements := parser.ParseStatementSequence()
	parser.closeBlock()
	return &ast.BlockStatement{
		Children: statements,
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}
