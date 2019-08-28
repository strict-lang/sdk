package parsing

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

// ParseIfStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) ParseIfStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.IfKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	condition, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err := parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.ParseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if !token.HasKeywordValue(parsing.token(), token.ElseKeyword) {
		return &ast.ConditionalStatement{
			Condition:    condition,
			Consequence:  body,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	parsing.advance()
	elseBody, err := parsing.parseElseIfOrBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &ast.ConditionalStatement{
		Condition:    condition,
		Consequence:  body,
		Alternative:  elseBody,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() (ast.Node, error) {
	if token.HasKeywordValue(parsing.token(), token.IfKeyword) {
		return parsing.ParseIfStatement(), nil
	}
	parsing.skipEndOfStatement()
	return parsing.ParseStatementBlock()
}

// ParseForStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parsing *Parsing) ParseForStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ForKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	initializerBeginToken := parsing.token()
	if token.IsIdentifierToken(initializerBeginToken) {
		if token.HasKeywordValue(parsing.peek(), token.FromKeyword) {
			return parsing.completeFromToStatement(beginOffset)
		}
	}
	return parsing.completeForEachStatement(beginOffset)
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops parsing.
func (parsing *Parsing) completeForEachStatement(beginOffset source.Offset) ast.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err := parsing.skipKeyword(token.InKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	value, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err := parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.ParseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &ast.ForEachLoopStatement{
		Field:        field,
		Enumeration:  value,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops parsing.
func (parsing *Parsing) completeFromToStatement(beginOffset source.Offset) ast.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err := parsing.skipKeyword(token.FromKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	from, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err := parsing.skipKeyword(token.ToKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	to, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err := parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.ParseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &ast.RangedLoopStatement{
		ValueField:   field,
		InitialValue: from,
		EndValue:     to,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// ParseYieldStatement parses a 'yield' statement. Yield statements add an
// element to an implicitly created list, which is returned by the method.
// Any kind of expression can be yielded.
func (parsing *Parsing) ParseYieldStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.YieldKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	rightHandSide, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &ast.YieldStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// ParseReturnStatement parses a 'return' statement. Return statements are
// part of the control flow and used to either return a value from a method
// or to end the method call within a branch, resulting in the remaining
// instructions to be ignored. The ReturnStatement is always the last statement
// within a StatementSequence / Branch.
func (parsing *Parsing) ParseReturnStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ReturnKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &ast.ReturnStatement{
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	rightHandSide, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &ast.ReturnStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() ast.Node {
	beginPosition := parsing.offset()
	method, err := parsing.ParseMethodDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginPosition, err)
	}
	return method
}

func (parsing *Parsing) ParseImportStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ImportKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	path := parsing.token()
	if !token.IsStringLiteralToken(path) {
		return parsing.createInvalidStatement(beginOffset, &UnexpectedTokenError{
			Expected: "Path",
			Token:    path,
		})
	}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &ast.ImportStatement{
			Path:         path.Value(),
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	parsing.advance()
	aliasOffset := parsing.offset()
	alias, err := parsing.parseImportAlias()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	aliasEnd := parsing.offset()
	parsing.skipEndOfStatement()
	return &ast.ImportStatement{
		Path: path.Value(),
		Alias: &ast.Identifier{
			Value: alias,
			NodePosition: &offsetPosition{
				begin: aliasOffset,
				end:   aliasEnd,
			},
		},
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseImportAlias() (string, error) {
	alias := parsing.token()
	if !token.IsIdentifierToken(alias) {
		return "", &UnexpectedTokenError{
			Expected: "Identifier",
			Token:    alias,
		}
	}
	parsing.advance()
	return alias.Value(), nil
}

var errNoAssign = errors.New("no assign")

func (parsing *Parsing) parseOptionalAssignValue() (ast.Node, error) {
	if !token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		return nil, errNoAssign
	}
	parsing.advance()
	return parsing.ParseExpression()
}

// keywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parsing *Parsing) keywordStatementParser(keyword token.Keyword) (func() ast.Node, bool) {
	switch keyword {
	case token.IfKeyword:
		return parsing.ParseIfStatement, true
	case token.ForKeyword:
		return parsing.ParseForStatement, true
	case token.YieldKeyword:
		return parsing.ParseYieldStatement, true
	case token.ReturnKeyword:
		return parsing.ParseReturnStatement, true
	case token.ImportKeyword:
		return parsing.ParseImportStatement, true
	case token.TestKeyword:
		return parsing.ParseTestStatement, true
	case token.AssertKeyword:
		return parsing.ParseAssertStatement, true
	case token.MethodKeyword:
		return parsing.parseNestedMethodDeclaration, true
	}
	return nil, false
}

// ParseKeywordStatement parses a statement that starts with a keyword.
func (parsing *Parsing) ParseKeywordStatement(keyword token.Keyword) ast.Node {
	function, ok := parsing.keywordStatementParser(keyword)
	if ok {
		return function()
	}
	parsing.reportError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "statement begin",
	})
	return &ast.InvalidStatement{}
}

func (parsing *Parsing) parseAssignStatement(operator token.Operator, leftHandSide ast.Node) (ast.Node, error) {
	beginOffset := parsing.offset()
	rightHandSide, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	parsing.skipEndOfStatement()
	return &ast.AssignStatement{
		Target:       leftHandSide,
		Value:        rightHandSide,
		Operator:     operator,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) ParseTestStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.TestKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	statements, err := parsing.ParseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &ast.TestStatement{
		NodePosition: parsing.createPosition(beginOffset),
		MethodName:   parsing.currentMethodName,
		Statements:   statements,
	}
}

func (parsing *Parsing) ParseAssertStatement() ast.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.AssertKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	expression, err := parsing.ParseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &ast.AssertStatement{
		NodePosition: parsing.createPosition(beginOffset),
		Expression:   expression,
	}
}

// ParseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) ParseInstructionStatement() (ast.Node, error) {
	leftHandSide, err := parsing.ParseExpression()
	if err != nil {
		return nil, err
	}
	switch operator := token.OperatorValue(parsing.token()); {
	case operator.IsAssign():
		parsing.skipEndOfStatement()
		return parsing.parseAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &ast.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token.DecrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &ast.DecrementStatement{Operand: leftHandSide}, nil
	}
	parsing.advance()
	return &ast.ExpressionStatement{
		Expression: leftHandSide,
	}, nil
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parsing *Parsing) ParseStatement() ast.Node {
	beginOffset := parsing.offset()
	switch current := parsing.token(); {
	case token.IsKeywordToken(current):
		return parsing.ParseKeywordStatement(token.KeywordValue(current))
	case token.IsIdentifierToken(current):
		fallthrough
	case token.IsOperatorToken(current):
		fallthrough
	case token.IsLiteralToken(current):
		statement, err := parsing.ParseInstructionStatement()
		if err != nil {
			return parsing.createInvalidStatement(beginOffset, err)
		}
		return statement
	default:
		return parsing.createInvalidStatement(beginOffset, ErrInvalidExpression)
	}
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parsing *Parsing) ParseStatementSequence() []ast.Node {
	var statements []ast.Node
	for {
		expectedIndent := parsing.block.Indent
		current := parsing.token()
		if token.IsEndOfFileToken(current) {
			break
		}
		if current.Indent() > expectedIndent {
			beginOffset := parsing.offset() - 1
			invalid := parsing.createInvalidStatement(beginOffset, &InvalidIndentationError{
				Token:    current,
				Expected: fmt.Sprintf("indent level of %d", expectedIndent),
			})
			statements = append(statements, invalid)
			break
		}
		if current.Indent() < expectedIndent {
			break
		}
		statement := parsing.ParseStatement()
		if _, ok := statement.(*ast.InvalidStatement); ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

// ParseStatementBlock parses a block of statements.
func (parsing *Parsing) ParseStatementBlock() (*ast.BlockStatement, error) {
	beginOffset := parsing.offset()
	indent := parsing.token().Indent()
	if indent < parsing.block.Indent {
		return nil, &InvalidIndentationError{
			Token:    parsing.token(),
			Expected: "indent bigger than 0",
		}
	}
	parsing.openBlock(indent)
	statements := parsing.ParseStatementSequence()
	parsing.closeBlock()
	return &ast.BlockStatement{
		Children:     statements,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}
