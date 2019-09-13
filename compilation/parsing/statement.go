package parsing

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.IfKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	condition, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	consequence, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if !token.HasKeywordValue(parsing.token(), token.ElseKeyword) {
		return &syntaxtree.ConditionalStatement{
			Condition:    condition,
			Consequence:  consequence,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	return parsing.parseConditionalStatementWithAlternative(
		beginOffset, condition, consequence)
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	beginOffset source.Offset, condition syntaxtree.Node, consequence syntaxtree.Node) syntaxtree.Node {

	parsing.advance()
	alternative, err := parsing.parseElseIfOrBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree.ConditionalStatement{
		Condition:    condition,
		Consequence:  consequence,
		Alternative:  alternative,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() (syntaxtree.Node, error) {
	if token.HasKeywordValue(parsing.token(), token.IfKeyword) {
		return parsing.parseConditionalStatement(), nil
	}
	parsing.skipEndOfStatement()
	return parsing.parseStatementBlock()
}

// parseForStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parsing *Parsing) parseForStatement() syntaxtree.Node {
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
func (parsing *Parsing) completeForEachStatement(beginOffset source.Offset) syntaxtree.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err = parsing.skipKeyword(token.InKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	value, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree.ForEachLoopStatement{
		Field:        field,
		Sequence:     value,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops parsing.
func (parsing *Parsing) completeFromToStatement(beginOffset source.Offset) syntaxtree.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err = parsing.skipKeyword(token.FromKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	from, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token.ToKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	to, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree.RangedLoopStatement{
		ValueField:   field,
		InitialValue: from,
		EndValue:     to,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseYieldStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.YieldKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.YieldStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseReturnStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ReturnKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &syntaxtree.ReturnStatement{
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.ReturnStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() syntaxtree.Node {
	beginPosition := parsing.offset()
	method, err := parsing.parseMethodDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginPosition, err)
	}
	return method
}

func (parsing *Parsing) parseImportStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ImportKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	switch next := parsing.token(); {
	case token.IsStringLiteralToken(next):
		return parsing.parseFileImport(beginOffset)
	case token.IsIdentifierToken(next):
		return parsing.parseIdentifierChainImport(beginOffset)
	default:
		return parsing.createInvalidStatement(beginOffset, &UnexpectedTokenError{
			Expected: "File or Class Path",
			Token:    next,
		})
	}
}

func (parsing *Parsing) parseIdentifierChainImport(beginOffset source.Offset) syntaxtree.Node {
	var chain []string
	for token.IsIdentifierToken(parsing.token()) {
		chain = append(chain, parsing.token().Value())
		parsing.advance()
		if token.HasOperatorValue(parsing.token(), token.DotOperator) {
			parsing.advance()
		}
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.ImportStatement{
		Target:       &syntaxtree.IdentifierChainImport{Chain: chain},
		Alias:        nil,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFileImport(beginOffset source.Offset) syntaxtree.Node {
	target := &syntaxtree.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &syntaxtree.ImportStatement{
			Target:       target,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	parsing.advance()
	return parsing.parseFileImportWithAlias(beginOffset, target)
}

func (parsing *Parsing) parseFileImportWithAlias(
	beginOffset source.Offset, target syntaxtree.ImportTarget) syntaxtree.Node {

	aliasOffset := parsing.offset()
	alias, err := parsing.parseImportAlias()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	aliasEnd := parsing.offset()
	parsing.skipEndOfStatement()
	return &syntaxtree.ImportStatement{
		Target: target,
		Alias: &syntaxtree.Identifier{
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

func (parsing *Parsing) parseOptionalAssignValue() (syntaxtree.Node, error) {
	if !token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		return nil, errNoAssign
	}
	parsing.advance()
	return parsing.parseExpression()
}

// findKeywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parsing *Parsing) findKeywordStatementParser(keyword token.Keyword) (func() syntaxtree.
	Node, bool) {
	switch keyword {
	case token.IfKeyword:
		return parsing.parseConditionalStatement, true
	case token.ForKeyword:
		return parsing.parseForStatement, true
	case token.YieldKeyword:
		return parsing.parseYieldStatement, true
	case token.ReturnKeyword:
		return parsing.parseReturnStatement, true
	case token.ImportKeyword:
		return parsing.parseImportStatement, true
	case token.TestKeyword:
		return parsing.parseTestStatement, true
	case token.AssertKeyword:
		return parsing.parseAssertStatement, true
	case token.MethodKeyword:
		return parsing.parseNestedMethodDeclaration, true
	case token.CreateKeyword:
		return parsing.maybeParseConstructorDeclaration()
	}
	return nil, false
}

// shouldParseConstructorDeclarations tells the parser whether it should
// attempt to parse a ConstructorDeclaration by determining if it is looking
// at a constructor call or declaration.
func (parsing *Parsing) shouldParseConstructorDeclaration() bool {
	return parsing.isAtBeginOfStatement && !token.IsIdentifierToken(parsing.peek())
}

func (parsing *Parsing) maybeParseConstructorDeclaration() (func() syntaxtree.Node, bool) {
	return parsing.parseConstructorDeclaration, parsing.shouldParseConstructorDeclaration()
}

func (parsing *Parsing) parseConstructorDeclaration() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.CreateKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parameters, err := parsing.parseParameterList()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree.ConstructorDeclaration{
		Parameters:   parameters,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// parseKeywordStatement parses a statement that starts with a keyword.
func (parsing *Parsing) parseKeywordStatement(keyword token.Keyword) syntaxtree.Node {
	function, ok := parsing.findKeywordStatementParser(keyword)
	if ok {
		return function()
	}
	parsing.reportError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "keyword that starts a statement",
	}, parsing.createTokenPosition())
	return &syntaxtree.InvalidStatement{}
}

func (parsing *Parsing) parseAssignStatement(
	operator token.Operator, leftHandSide syntaxtree.Node) (syntaxtree.Node, error) {

	beginOffset := parsing.offset()
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.AssignStatement{
		Target:       leftHandSide,
		Value:        rightHandSide,
		Operator:     operator,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseTestStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.TestKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	statements, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree.TestStatement{
		NodePosition: parsing.createPosition(beginOffset),
		MethodName:   parsing.currentMethodName,
		Statements:   statements,
	}
}

func (parsing *Parsing) parseAssertStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.AssertKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	expression, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.AssertStatement{
		NodePosition: parsing.createPosition(beginOffset),
		Expression:   expression,
	}
}

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() (syntaxtree.Node, error) {
	leftHandSide, err := parsing.parseExpression()
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
		return &syntaxtree.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token.DecrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &syntaxtree.DecrementStatement{Operand: leftHandSide}, nil
	}
	parsing.advance()
	return &syntaxtree.ExpressionStatement{
		Expression: leftHandSide,
	}, nil
}

func (parsing *Parsing) isKeywordStatementToken(entry token.Token) bool {
	if !token.IsKeywordToken(entry) {
		return false
	}
	if token.HasKeywordValue(entry, token.CreateKeyword) {
		return parsing.shouldParseConstructorDeclaration()
	}
	return true
}

func (parsing *Parsing) isKeywordExpressionToken(entry token.Token) bool {
	if !token.IsKeywordToken(entry) {
		return false
	}
	if token.HasKeywordValue(entry, token.CreateKeyword) {
		return !parsing.shouldParseConstructorDeclaration()
	}
	return true
}

func (parsing *Parsing) shouldParseFieldDeclaration() bool {
	return parsing.isAtBeginOfStatement && parsing.couldBeLookingAtTypeName()
}

func (parsing *Parsing) parseFieldDeclarationOrDefinition() syntaxtree.Node {
	beginOffset := parsing.offset()
	declaration, err := parsing.parseFieldDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if !token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		parsing.skipEndOfStatement()
		return declaration
	}
	parsing.advance()
	return parsing.completeParsingFieldDefinition(beginOffset, declaration)
}
func (parsing *Parsing) completeParsingFieldDefinition(
	beginOffset source.Offset, declaration syntaxtree.Node) syntaxtree.Node {

	value, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree.AssignStatement{
		Target:       declaration,
		Value:        value,
		Operator:     token.AssignOperator,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFieldDeclaration() (syntaxtree.Node, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	fieldName, err := parsing.expectAnyIdentifier()
	if err != nil {
		return nil, err
	}
	parsing.advance()
	return &syntaxtree.FieldDeclaration{
		Name:         fieldName,
		TypeName:     typeName,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parsing *Parsing) parseStatement() syntaxtree.Node {
	beginOffset := parsing.offset()
	switch current := parsing.token(); {
	case parsing.isKeywordStatementToken(current):
		return parsing.parseKeywordStatement(token.KeywordValue(current))
	case parsing.isKeywordExpressionToken(current):
		fallthrough
	case token.IsIdentifierToken(current):
		if parsing.shouldParseFieldDeclaration() {
			return parsing.parseFieldDeclarationOrDefinition()
		}
		fallthrough
	case token.IsOperatorToken(current):
		fallthrough
	case token.IsLiteralToken(current):
		statement, err := parsing.parseInstructionStatement()
		if err != nil {
			return parsing.createInvalidStatement(beginOffset, err)
		}
		return statement
	default:
		err := fmt.Errorf("expected begin of statement or expression but got: %s", current)
		return parsing.createInvalidStatement(beginOffset, err)
	}
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parsing *Parsing) parseStatementSequence() []syntaxtree.Node {
	var statements []syntaxtree.Node
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
		statement := parsing.parseStatement()
		if _, ok := statement.(*syntaxtree.InvalidStatement); ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

// ParseStatementBlock parses a block of statements.
func (parsing *Parsing) parseStatementBlock() (*syntaxtree.BlockStatement, error) {
	beginOffset := parsing.offset()
	indent := parsing.token().Indent()
	if indent < parsing.block.Indent {
		return nil, &InvalidIndentationError{
			Token:    parsing.token(),
			Expected: "indent bigger than 0",
		}
	}
	parsing.openBlock(indent)
	statements := parsing.parseStatementSequence()
	parsing.closeBlock()
	return &syntaxtree.BlockStatement{
		Children:     statements,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}
