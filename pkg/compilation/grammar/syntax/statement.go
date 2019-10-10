package syntax

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/syntax/tree"
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() tree.Node {
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
		return &tree.ConditionalStatement{
			Condition:    condition,
			Consequence:  consequence,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	return parsing.parseConditionalStatementWithAlternative(
		beginOffset, condition, consequence)
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	beginOffset input.Offset, condition tree.Node, consequence tree.Node) tree.Node {

	parsing.advance()
	alternative, err := parsing.parseElseIfOrBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &tree.ConditionalStatement{
		Condition:    condition,
		Consequence:  consequence,
		Alternative:  alternative,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() (tree.Node, error) {
	if token.HasKeywordValue(parsing.token(), token.IfKeyword) {
		return parsing.parseConditionalStatement(), nil
	}
	parsing.skipEndOfStatement()
	return parsing.parseStatementBlock()
}

// parseForStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parsing *Parsing) parseForStatement() tree.Node {
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
// This method completes the loops grammar.
func (parsing *Parsing) completeForEachStatement(beginOffset input.Offset) tree.Node {
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
	return &tree.ForEachLoopStatement{
		Field:        field,
		Sequence:     value,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops grammar.
func (parsing *Parsing) completeFromToStatement(beginOffset input.Offset) tree.Node {
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
	return &tree.RangedLoopStatement{
		ValueField:   field,
		InitialValue: from,
		EndValue:     to,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseYieldStatement() tree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.YieldKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &tree.YieldStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseReturnStatement() tree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.ReturnKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &tree.ReturnStatement{
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &tree.ReturnStatement{
		Value:        rightHandSide,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() tree.Node {
	beginPosition := parsing.offset()
	method, err := parsing.parseMethodDeclaration()
	if err != nil {
		return parsing.createInvalidStatement(beginPosition, err)
	}
	return method
}

func (parsing *Parsing) parseImportStatement() tree.Node {
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

func (parsing *Parsing) parseIdentifierChainImport(beginOffset input.Offset) tree.Node {
	var chain []string
	for token.IsIdentifierToken(parsing.token()) {
		chain = append(chain, parsing.token().Value())
		parsing.advance()
		if token.HasOperatorValue(parsing.token(), token.DotOperator) {
			parsing.advance()
		}
	}
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target:       &tree.IdentifierChainImport{Chain: chain},
		Alias:        nil,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFileImport(beginOffset input.Offset) tree.Node {
	target := &tree.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &tree.ImportStatement{
			Target:       target,
			NodePosition: parsing.createPosition(beginOffset),
		}
	}
	parsing.advance()
	return parsing.parseFileImportWithAlias(beginOffset, target)
}

func (parsing *Parsing) parseFileImportWithAlias(
	beginOffset input.Offset, target tree.ImportTarget) tree.Node {

	aliasOffset := parsing.offset()
	alias, err := parsing.parseImportAlias()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	aliasEnd := parsing.offset()
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target: target,
		Alias: &tree.Identifier{
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

func (parsing *Parsing) parseOptionalAssignValue() (tree.Node, error) {
	if !token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		return nil, errNoAssign
	}
	parsing.advance()
	return parsing.parseExpression()
}

// findKeywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parsing *Parsing) findKeywordStatementParser(keyword token.Keyword) (func() tree.Node, bool) {
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
	return !token.IsIdentifierToken(parsing.peek())
}

func (parsing *Parsing) maybeParseConstructorDeclaration() (func() tree.Node, bool) {
	return parsing.parseConstructorDeclaration, parsing.shouldParseConstructorDeclaration()
}

func (parsing *Parsing) parseConstructorDeclaration() tree.Node {
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
	return &tree.ConstructorDeclaration{
		Parameters:   parameters,
		Body:         body,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

// parseKeywordStatement parses a statement that starts with a keyword.
func (parsing *Parsing) parseKeywordStatement(keyword token.Keyword) tree.Node {
	function, ok := parsing.findKeywordStatementParser(keyword)
	if ok {
		return function()
	}
	parsing.reportError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "keyword that starts a statement",
	}, parsing.createTokenPosition())
	return &tree.InvalidStatement{}
}

func (parsing *Parsing) parseAssignStatement(
	operator token.Operator, leftHandSide tree.Node) (tree.Node, error) {

	beginOffset := parsing.offset()
	rightHandSide, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err), err
	}
	parsing.skipEndOfStatement()
	return &tree.AssignStatement{
		Target:       leftHandSide,
		Value:        rightHandSide,
		Operator:     operator,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseTestStatement() tree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.TestKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	statements, err := parsing.parseStatementBlock()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	return &tree.TestStatement{
		NodePosition: parsing.createPosition(beginOffset),
		MethodName:   parsing.currentMethodName,
		Statements:   statements,
	}
}

func (parsing *Parsing) parseAssertStatement() tree.Node {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.AssertKeyword); err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	expression, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &tree.AssertStatement{
		NodePosition: parsing.createPosition(beginOffset),
		Expression:   expression,
	}
}

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() (tree.Node, error) {
	leftHandSide, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	switch operator := token.OperatorValue(parsing.token()); {
	case operator.IsAssign():
		parsing.advance() // ?
		return parsing.parseAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &tree.IncrementStatement{Operand: leftHandSide}, nil
	case operator == token.DecrementOperator:
		parsing.advance()
		parsing.skipEndOfStatement()
		return &tree.DecrementStatement{Operand: leftHandSide}, nil
	}
	parsing.advance()
	return &tree.ExpressionStatement{
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
	return /* grammar.isAtBeginOfStatement && */ parsing.couldBeLookingAtTypeName()
}

func (parsing *Parsing) parseFieldDeclarationOrDefinition() tree.Node {
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
	beginOffset input.Offset, declaration tree.Node) tree.Node {

	value, err := parsing.parseExpression()
	if err != nil {
		return parsing.createInvalidStatement(beginOffset, err)
	}
	parsing.skipEndOfStatement()
	return &tree.AssignStatement{
		Target:       declaration,
		Value:        value,
		Operator:     token.AssignOperator,
		NodePosition: parsing.createPosition(beginOffset),
	}
}

func (parsing *Parsing) parseFieldDeclaration() (tree.Node, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	return parsing.parseFieldDeclarationWithTypeName(beginOffset, typeName)
}

func (parsing *Parsing) parseFieldDeclarationWithTypeName(
	beginOffset input.Offset, typeName tree.TypeName) (tree.Node, error) {

	fieldName, err := parsing.expectAnyIdentifier()
	if err != nil {
		return nil, err
	}
	parsing.advance()
	declaration := &tree.FieldDeclaration{
		Name:         fieldName,
		TypeName:     typeName,
		NodePosition: parsing.createPosition(beginOffset),
	}
	if token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		return parsing.parseFieldDefinition(beginOffset, declaration)
	}
	return declaration, nil
}

func (parsing *Parsing) parseFieldDefinition(
	beginOffset input.Offset, declaration *tree.FieldDeclaration) (tree.Node, error) {

	if err := parsing.skipOperator(token.AssignOperator); err != nil {
		return nil, err
	}
	value, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	return &tree.AssignStatement{
		Target:       declaration,
		Value:        value,
		Operator:     token.AssignOperator,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parsing *Parsing) parseStatement() tree.Node {
	beginOffset := parsing.offset()
	switch current := parsing.token(); {
	case parsing.isKeywordStatementToken(current):
		return parsing.parseKeywordStatement(token.KeywordValue(current))
	case parsing.isKeywordExpressionToken(current):
		fallthrough
	case token.IsIdentifierToken(current):
		if parsing.shouldParseFieldDeclaration() {
			node, err := parsing.parseFieldDeclarationOrListAccess()
			if err != nil {
				return parsing.createInvalidStatement(beginOffset, err)
			}
			return node
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

// parseFieldDeclarationOrListAccess either parses a FieldDeclaration or a ListSelect.
// Both constructs are similar when only being able to peek one token.
func (parsing *Parsing) parseFieldDeclarationOrListAccess() (tree.Node, error) {
	beginOffset := parsing.offset()
	baseTypeOrAccessedField := parsing.token()
	parsing.advance()
	if (!token.HasOperatorValue(parsing.token(), token.LeftBracketOperator)) ||
		(token.OperatorValue(parsing.peek()) == token.RightBracketOperator) {

		return parsing.parseFieldDeclarationFromBaseTypeName(beginOffset, baseTypeOrAccessedField)
	}
	// This method has to be used because of selector chaining and calls.
	operation, err := parsing.parseOperationsOnOperand(&tree.Identifier{
		Value:        baseTypeOrAccessedField.Value(),
		NodePosition: parsing.createPosition(beginOffset),
	})
	if err != nil {
		return nil, err
	}
	return parsing.parseOperationOrAssign(operation)
}

func (parsing *Parsing) parseFieldDeclarationFromBaseTypeName(
	beginOffset input.Offset, baseTypeName token.Token) (tree.Node, error){

	typeName, err := parsing.parseTypeNameFromBaseIdentifier(beginOffset, baseTypeName)
	if err != nil {
		return nil, err
	}
	return parsing.parseFieldDeclarationWithTypeName(beginOffset, typeName)
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parsing *Parsing) parseStatementSequence() []tree.Node {
	var statements []tree.Node
	for {
		expectedIndent := parsing.block.Indent
		current := parsing.token()
		if token.IsEndOfFileToken(current) {
			break
		}
		// FIXME
		if token.IsEndOfStatementToken(current) {
			parsing.advance()
			continue
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
		if _, ok := statement.(*tree.InvalidStatement); ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

// ParseStatementBlock parses a block of statements.
func (parsing *Parsing) parseStatementBlock() (*tree.BlockStatement, error) {
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
	return &tree.BlockStatement{
		Children:     statements,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}
