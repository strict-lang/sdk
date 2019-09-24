package parsing

import (
	"errors"
	"fmt"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.IfKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	condition, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token2.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	consequence, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if !token2.HasKeywordValue(parsing.token(), token2.ElseKeyword) {
		return &syntaxtree2.ConditionalStatement{
			Condition:    condition,
			Consequence:  consequence,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	return parsing.parseConditionalStatementWithAlternative(
		beginOffset, condition, consequence)
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	beginOffset source2.Offset, condition syntaxtree2.Node, consequence syntaxtree2.Node) syntaxtree2.Node {

	parsing.advance()
	alternative, err := parsing.parseElseIfOrBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree2.ConditionalStatement{
		Condition:    condition,
		Consequence:  consequence,
		Alternative:  alternative,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() (syntaxtree2.Node, error) {
	if token2.HasKeywordValue(parsing.token(), token2.IfKeyword) {
		return parsing.parseConditionalStatement(), nil
	}
	parsing.skipEndOfStatement()
	return parsing.parseStatementBlock()
}

// parseForStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parsing *Parsing) parseForStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.ForKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	initializerBeginToken := parsing.token()
	if token2.IsIdentifierToken(initializerBeginToken) {
		if token2.HasKeywordValue(parsing.peek(), token2.FromKeyword) {
			return parsing.completeFromToStatement(beginOffset)
		}
	}
	return parsing.completeForEachStatement(beginOffset)
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops parsing.
func (parsing *Parsing) completeForEachStatement(beginOffset source2.Offset) syntaxtree2.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err = parsing.skipKeyword(token2.InKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	value, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token2.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree2.ForEachLoopStatement{
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
func (parsing *Parsing) completeFromToStatement(beginOffset source2.Offset) syntaxtree2.Node {
	field, err := parsing.expectAnyIdentifier()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.advance()
	if err = parsing.skipKeyword(token2.FromKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	from, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token2.ToKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	to, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if err = parsing.skipKeyword(token2.DoKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	body, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree2.RangedLoopStatement{
		ValueField:   field,
		InitialValue: from,
		EndValue:     to,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseYieldStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.YieldKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.YieldStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseReturnStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.ReturnKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if token2.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &syntaxtree2.ReturnStatement{
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.ReturnStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() syntaxtree2.Node {
	beginPosition := parsing.offset()
	method, err := parsing.parseMethodDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginPosition, err)
	}
	return method
}

func (parsing *Parsing) parseImportStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.ImportKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	switch next := parsing.token(); {
	case token2.IsStringLiteralToken(next):
		return parsing.parseFileImport(beginOffset)
	case token2.IsIdentifierToken(next):
		return parsing.parseIdentifierChainImport(beginOffset)
	default:
		return parsing.createInvalidStatement(beginOffset, &UnexpectedTokenError{
			Expected: "File or Class Path",
			Token:    next,
		})
	}
}

func (parsing *Parsing) parseIdentifierChainImport(beginOffset source2.Offset) syntaxtree2.Node {
	var chain []string
	for token2.IsIdentifierToken(parsing.token()) {
		chain = append(chain, parsing.token().Value())
		parsing.advance()
		if token2.HasOperatorValue(parsing.token(), token2.DotOperator) {
			parsing.advance()
		}
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.ImportStatement{
		Target:       &syntaxtree2.IdentifierChainImport{Chain: chain},
		Alias:        nil,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFileImport(beginOffset source2.Offset) syntaxtree2.Node {
	target := &syntaxtree2.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token2.HasKeywordValue(parsing.token(), token2.AsKeyword) {
		parsing.skipEndOfStatement()
		return &syntaxtree2.ImportStatement{
			Target:       target,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	parsing.advance()
	return parsing.parseFileImportWithAlias(beginOffset, target)
}

func (parsing *Parsing) parseFileImportWithAlias(
	beginOffset source2.Offset, target syntaxtree2.ImportTarget) syntaxtree2.Node {

	aliasOffset := parsing.offset()
	alias, err := parsing.parseImportAlias()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	aliasEnd := parsing.offset()
	parsing.skipEndOfStatement()
	return &syntaxtree2.ImportStatement{
		Target: target,
		Alias: &syntaxtree2.Identifier{
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
	if !token2.IsIdentifierToken(alias) {
		return "", &UnexpectedTokenError{
			Expected: "Identifier",
			Token:    alias,
		}
	}
	parsing.advance()
	return alias.Value(), nil
}

var errNoAssign = errors.New("no assign")

func (parsing *Parsing) parseOptionalAssignValue() (syntaxtree2.Node, error) {
	if !token2.HasOperatorValue(parsing.token(), token2.AssignOperator) {
		return nil, errNoAssign
	}
	parsing.advance()
	return parsing.parseExpression()
}

// findKeywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parsing *Parsing) findKeywordStatementParser(keyword token2.Keyword) (func() syntaxtree2.Node, bool) {
	switch keyword {
	case token2.IfKeyword:
		return parsing.parseConditionalStatement, true
	case token2.ForKeyword:
		return parsing.parseForStatement, true
	case token2.YieldKeyword:
		return parsing.parseYieldStatement, true
	case token2.ReturnKeyword:
		return parsing.parseReturnStatement, true
	case token2.ImportKeyword:
		return parsing.parseImportStatement, true
	case token2.TestKeyword:
		return parsing.parseTestStatement, true
	case token2.AssertKeyword:
		return parsing.parseAssertStatement, true
	case token2.MethodKeyword:
		return parsing.parseNestedMethodDeclaration, true
	case token2.CreateKeyword:
		return parsing.maybeParseConstructorDeclaration()
	}
	return nil, false
}

// shouldParseConstructorDeclarations tells the parser whether it should
// attempt to parse a ConstructorDeclaration by determining if it is looking
// at a constructor call or declaration.
func (parsing *Parsing) shouldParseConstructorDeclaration() bool {
	return parsing.isAtBeginOfStatement && !token2.IsIdentifierToken(parsing.peek())
}

func (parsing *Parsing) maybeParseConstructorDeclaration() (func() syntaxtree2.Node, bool) {
	return parsing.parseConstructorDeclaration, parsing.shouldParseConstructorDeclaration()
}

func (parsing *Parsing) parseConstructorDeclaration() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.CreateKeyword); err != nil {
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
	return &syntaxtree2.ConstructorDeclaration{
		Parameters:   parameters,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// parseKeywordStatement parses a statement that starts with a keyword.
func (parsing *Parsing) parseKeywordStatement(keyword token2.Keyword) syntaxtree2.Node {
	function, ok := parsing.findKeywordStatementParser(keyword)
	if ok {
		return function()
	}
	parsing.reportError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "keyword that starts a statement",
	}, parsing.createTokenPosition())
	return &syntaxtree2.InvalidStatement{}
}

func (parsing *Parsing) parseAssignStatement(
	operator token2.Operator, leftHandSide syntaxtree2.Node) (syntaxtree2.Node, error) {

	beginOffset := parsing.offset()
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.AssignStatement{
		Target:       leftHandSide,
		Value:        rightHandSide,
		Operator:     operator,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseTestStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.TestKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	statements, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &syntaxtree2.TestStatement{
		NodePosition: parsing.createPosition(beginOffset),
		MethodName:   parsing.currentMethodName,
		Statements:   statements,
	}
}

func (parsing *Parsing) parseAssertStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.AssertKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	expression, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.AssertStatement{
		NodePosition: parsing.createPosition(beginOffset),
		Expression:   expression,
	}
}

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() (syntaxtree2.Node, error) {
	leftHandSide, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	switch operator := token2.OperatorValue(parsing.token()); {
	case operator.IsAssign():
		parsing.skipEndOfStatement()
		return parsing.parseAssignStatement(operator, leftHandSide)
	case operator == token2.IncrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &syntaxtree2.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token2.DecrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &syntaxtree2.DecrementStatement{Operand: leftHandSide}, nil
	}
	parsing.advance()
	return &syntaxtree2.ExpressionStatement{
		Expression: leftHandSide,
	}, nil
}

func (parsing *Parsing) isKeywordStatementToken(entry token2.Token) bool {
	if !token2.IsKeywordToken(entry) {
		return false
	}
	if token2.HasKeywordValue(entry, token2.CreateKeyword) {
		return parsing.shouldParseConstructorDeclaration()
	}
	return true
}

func (parsing *Parsing) isKeywordExpressionToken(entry token2.Token) bool {
	if !token2.IsKeywordToken(entry) {
		return false
	}
	if token2.HasKeywordValue(entry, token2.CreateKeyword) {
		return !parsing.shouldParseConstructorDeclaration()
	}
	return true
}

func (parsing *Parsing) shouldParseFieldDeclaration() bool {
	return parsing.isAtBeginOfStatement && parsing.couldBeLookingAtTypeName()
}

func (parsing *Parsing) parseFieldDeclarationOrDefinition() syntaxtree2.Node {
	beginOffset := parsing.offset()
	declaration, err := parsing.parseFieldDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if !token2.HasOperatorValue(parsing.token(), token2.AssignOperator) {
		parsing.skipEndOfStatement()
		return declaration
	}
	parsing.advance()
	return parsing.completeParsingFieldDefinition(beginOffset, declaration)
}
func (parsing *Parsing) completeParsingFieldDefinition(
	beginOffset source2.Offset, declaration syntaxtree2.Node) syntaxtree2.Node {

	value, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &syntaxtree2.AssignStatement{
		Target:       declaration,
		Value:        value,
		Operator:     token2.AssignOperator,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFieldDeclaration() (syntaxtree2.Node, error) {
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
	return &syntaxtree2.FieldDeclaration{
		Name:         fieldName,
		TypeName:     typeName,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parsing *Parsing) parseStatement() syntaxtree2.Node {
	beginOffset := parsing.offset()
	switch current := parsing.token(); {
	case parsing.isKeywordStatementToken(current):
		return parsing.parseKeywordStatement(token2.KeywordValue(current))
	case parsing.isKeywordExpressionToken(current):
		fallthrough
	case token2.IsIdentifierToken(current):
		if parsing.shouldParseFieldDeclaration() {
			return parsing.parseFieldDeclarationOrDefinition()
		}
		fallthrough
	case token2.IsOperatorToken(current):
		fallthrough
	case token2.IsLiteralToken(current):
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
func (parsing *Parsing) parseStatementSequence() []syntaxtree2.Node {
	var statements []syntaxtree2.Node
	for {
		expectedIndent := parsing.block.Indent
		current := parsing.token()
		if token2.IsEndOfFileToken(current) {
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
		if _, ok := statement.(*syntaxtree2.InvalidStatement); ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

// ParseStatementBlock parses a block of statements.
func (parsing *Parsing) parseStatementBlock() (*syntaxtree2.BlockStatement, error) {
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
	return &syntaxtree2.BlockStatement{
		Children:     statements,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}
