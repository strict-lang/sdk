package parser

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// ParseIfStatement parses a conditional statement and it's optional else-clause.
func (parser *Parser) ParseIfStatement() ast.Node {
	if err := parser.skipKeyword(token.IfKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	condition, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if !token.HasKeywordValue(parser.token(), token.ElseKeyword) {
		return &ast.ConditionalStatement{
			Condition: condition,
			Body:      body,
		}
	}
	parser.advance()
	elseBody, err := parser.parseElseIfOrBlock()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return &ast.ConditionalStatement{
		Condition: condition,
		Body:      body,
		Else:      elseBody,
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
	if err := parser.skipKeyword(token.ForKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	initializerBeginToken := parser.token()
	if token.IsIdentifierToken(initializerBeginToken) {
		if token.HasKeywordValue(parser.peek(), token.FromKeyword) {
			return parser.completeFromToStatement()
		}
	}
	return parser.completeForEachStatement()
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops parsing.
func (parser *Parser) completeForEachStatement() ast.Node {
	field, err := parser.expectAnyIdentifier()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.advance()
	if err := parser.skipKeyword(token.InKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	value, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if err := parser.skipKeyword(token.DoKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return &ast.ForeachLoopStatement{
		Field:  field,
		Target: value,
		Body:   body,
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops parsing.
func (parser *Parser) completeFromToStatement() ast.Node {
	field, err := parser.expectAnyIdentifier()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.advance()
	if err := parser.skipKeyword(token.FromKeyword); err != nil {
		return parser.createInvalidStatement(err)
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
	parser.skipEndOfStatement()
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return &ast.FromToLoopStatement{
		Field: field,
		From:  from,
		To:    to,
		Body:  body,
	}
}

// ParseYieldStatement parses a 'yield' statement. Yield statements add an
// element to an implicitly created list, which is returned by the method.
// Any kind of expression can be yielded.
func (parser *Parser) ParseYieldStatement() ast.Node {
	if err := parser.skipKeyword(token.YieldKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
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
	if err := parser.skipKeyword(token.ReturnKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	if token.IsEndOfStatementToken(parser.token()) {
		parser.advance()
		return &ast.ReturnStatement{}
	}
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
	return &ast.ReturnStatement{
		Value: rightHandSide,
	}
}

func (parser *Parser) parseNestedMethod() ast.Node {
	method, err := parser.ParseMethodDeclaration()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	return method
}

func (parser *Parser) ParseImportStatement() ast.Node {
	if err := parser.skipKeyword(token.ImportKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	path := parser.token()
	if !token.IsStringLiteralToken(path) {
		return parser.createInvalidStatement(&UnexpectedTokenError{
			Expected: "Path",
			Token: path,
		})
	}
	parser.advance()
	if !token.HasKeywordValue(parser.token(), token.AsKeyword) {
		parser.skipEndOfStatement()
		return &ast.ImportStatement{
			Path: path.Value(),
		}
	}
	parser.advance()
	alias, err := parser.parseImportAlias()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
	return &ast.ImportStatement{
		Path: path.Value(),
		Alias: ast.NewIdentifier(alias),
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

func (parser *Parser) ParseSharedVariableDeclaration() ast.Node {
	if err := parser.skipKeyword(token.SharedKeyword); err != nil {
		return parser.createInvalidStatement(err)
	}
	typeName, err := parser.ParseTypeName()
	if err != nil {
		return parser.createInvalidStatement(err)
	}
	if !token.IsIdentifierToken(parser.token()) {
		return parser.createInvalidStatement(&UnexpectedTokenError{
			Expected: "identifier",
			Token: parser.token(),
		})
	}
	variableName := parser.token().Value()
	parser.advance()
	assignedValue, err := parser.parseOptionalAssignValue()
	if err != nil && err != errNoAssign {
		return parser.createInvalidStatement(err)
	}
	parser.skipEndOfStatement()
	return &ast.SharedVariableDeclaration{
		Type: typeName,
		Name: ast.NewIdentifier(variableName),
		InitialValue: assignedValue,
	}
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
	case token.MethodKeyword:
		return parser.parseNestedMethod, true
	case token.IfKeyword:
		return parser.ParseIfStatement, true
	case token.ForKeyword:
		return parser.ParseForStatement, true
	case token.YieldKeyword:
		return parser.ParseYieldStatement, true
	case token.ReturnKeyword:
		return parser.ParseReturnStatement, true
	case token.ImportKeyword:
		return parser.ParseImportStatement, true
	case token.SharedKeyword:
		return parser.ParseSharedVariableDeclaration, true
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
	rightHandSide, err := parser.ParseExpression()
	if err != nil {
		return &ast.InvalidStatement{}, err
	}
	parser.skipEndOfStatement()
	return &ast.AssignStatement{
		Target:   leftHandSide,
		Value:    rightHandSide,
		Operator: operator,
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
			return parser.createInvalidStatement(err)
		}
		// if !token.IsEndOfStatementToken(parser.token()) {
		// }
		return statement
	default:
		return &ast.InvalidStatement{}
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
			invalid := parser.createInvalidStatement(&InvalidIndentationError{
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
	}, nil
}
