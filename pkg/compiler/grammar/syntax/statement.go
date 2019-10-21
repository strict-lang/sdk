package syntax

import (
	"errors"
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() *tree.ConditionalStatement {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.IfKeyword)
	condition := parsing.parseExpression()
	parsing.skipKeyword(token.DoKeyword)
	parsing.skipEndOfStatement()
	consequence := parsing.parseStatementBlock()
	return parsing.parseElseClauseIfPresent(beginOffset, condition, consequence)
}

func (parsing *Parsing) parseElseClauseIfPresent(
	beginOffset input.Offset,
	condition tree.Expression,
	consequence tree.Statement) *tree.ConditionalStatement {

	if token.HasKeywordValue(parsing.token(), token.ElseKeyword) {
		return parsing.parseConditionalStatementWithAlternative(
			beginOffset, condition, consequence)
	}
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Region:      parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	beginOffset input.Offset,
	condition tree.Node,
	consequence tree.Node) *tree.ConditionalStatement {

	parsing.advance()
	alternative := parsing.parseElseIfOrBlock()
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
		Region:      parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() tree.Node {
	if token.HasKeywordValue(parsing.token(), token.IfKeyword) {
		return parsing.parseConditionalStatement()
	}
	parsing.skipEndOfStatement()
	return parsing.parseStatementBlock()
}

// parseLoopStatement parses a loop statement, which starts with the
// ForKeyword. The statement may either be a FromToLoopStatement or
// a ForEachLoopStatement.
func (parsing *Parsing) parseLoopStatement() tree.Node {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.ForKeyword)
	if parsing.isLookingAtRangedLoop() {
		return parsing.completeFromToStatement(beginOffset)
	}
	return parsing.completeForEachStatement(beginOffset)
}

func (parsing* Parsing) isLookingAtRangedLoop() bool {
	return token.IsIdentifierToken(parsing.token()) &&
		token.HasKeywordValue(parsing.peek(), token.FromKeyword)
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops grammar.
func (parsing *Parsing) completeForEachStatement(beginOffset input.Offset) *tree.ForEachLoopStatement{
	field := parsing.parseIdentifier()
	parsing.skipKeyword(token.InKeyword)
	value := parsing.parseExpression()
	parsing.skipKeyword(token.DoKeyword)
	parsing.skipEndOfStatement()
	body := parsing.parseStatementBlock()
	return &tree.ForEachLoopStatement{
		Field:    field,
		Sequence: value,
		Body:     body,
		Region:   parsing.createRegion(beginOffset),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops grammar.
func (parsing *Parsing) completeFromToStatement(beginOffset input.Offset) *tree.RangedLoopStatement {
	field := parsing.parseIdentifier()
	parsing.skipKeyword(token.FromKeyword)
	begin := parsing.parseExpression()
	parsing.skipKeyword(token.ToKeyword)
	end := parsing.parseExpression()
	parsing.skipKeyword(token.DoKeyword)
	parsing.skipEndOfStatement()
	body := parsing.parseStatementBlock()
	return &tree.RangedLoopStatement{
		Field:  field,
		Begin:  begin,
		End:    end,
		Body:   body,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseYieldStatement() *tree.YieldStatement {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.YieldKeyword)
	rightHandSide := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.YieldStatement{
		Value:  rightHandSide,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseReturnStatement() *tree.ReturnStatement {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.ReturnKeyword)
	defer parsing.skipEndOfStatement()
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &tree.ReturnStatement{
			Region: parsing.createRegion(beginOffset),
		}
	}
	rightHandSide := parsing.parseExpression()
	return &tree.ReturnStatement{
		Value:  rightHandSide,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() *tree.MethodDeclaration {
	return parsing.parseMethodDeclaration()
}

func (parsing *Parsing) parseImportStatement() *tree.ImportStatement {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.ImportKeyword)
	switch next := parsing.token(); {
	case token.IsStringLiteralToken(next):
		return parsing.parseFileImport(beginOffset)
	case token.IsIdentifierToken(next):
		return parsing.parseIdentifierChainImport(beginOffset)
	default:
		parsing.throwError(&UnexpectedTokenError{
			Expected: "File or Class Path",
			Token:    next,
		})
		return nil
	}
}

func (parsing *Parsing) parseIdentifierChainImport(beginOffset input.Offset) *tree.ImportStatement {
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
		Target: &tree.IdentifierChainImport{Chain: chain},
		Alias:  nil,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseFileImport(beginOffset input.Offset) *tree.ImportStatement {
	target := &tree.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &tree.ImportStatement{
			Target: target,
			Region: parsing.createRegion(beginOffset),
		}
	}
	parsing.advance()
	return parsing.parseFileImportWithAlias(beginOffset, target)
}

func (parsing *Parsing) parseFileImportWithAlias(
	beginOffset input.Offset, target tree.ImportTarget) *tree.ImportStatement {

	alias := parsing.parseImportAlias()
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target: target,
		Alias: alias,
		Region: parsing.createRegion(beginOffset),
	}
}

func (parsing *Parsing) parseImportAlias() *tree.Identifier {
	return parsing.parseIdentifier()
}

var errNoAssign = errors.New("no assign")

func (parsing *Parsing) parseOptionalAssignValue() tree.Expression {
	parsing.skipOperator(token.AssignOperator)
	return parsing.parseExpression()
}

type keywordStatementParser func(*Parsing) tree.Node

var keywordStatementParserTable = map[token.Keyword] keywordStatementParser {
	token.IfKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseConditionalStatement()
	},
	token.ForKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseLoopStatement()
	},
	token.YieldKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseYieldStatement()
	},
	token.ReturnKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseReturnStatement()
	},
	token.ImportKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseImportStatement()
	},
	token.AssertKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseAssertStatement()
	},
	token.TestKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseTestStatement()
	},
	token.MethodKeyword: func(parsing *Parsing) tree.Node {
		return parsing.parseMethodDeclaration()
	},
}

// findKeywordStatementParser returns a function that parses statements based on a passed
// keyword. Most of the keywords start a statement. The returned bool is true, if a
// function has been found.
func (parsing *Parsing) findKeywordStatementParser(
	keyword token.Keyword) (keywordStatementParser, bool) {
	if keyword == token.CreateKeyword {
		return parsing.maybeParseConstructorDeclaration()
	}
	parser, found := keywordStatementParserTable[keyword]
	return parser, found
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
		Parameters: parameters,
		Child:      body,
		Region:     parsing.createRegion(beginOffset),
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
	}, parsing.createRegionFromToken())
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
		Target:   leftHandSide,
		Value:    rightHandSide,
		Operator: operator,
		Region:   parsing.createRegion(beginOffset),
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
		Region:     parsing.createRegion(beginOffset),
		MethodName: parsing.currentMethodName,
		Child:      statements,
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
		Region:     parsing.createRegion(beginOffset),
		Expression: expression,
	}
}

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() (tree.Node, error) {
	beginOffset := parsing.offset()
	leftHandSide, err := parsing.parseExpression()
	if err != nil {
		return nil, err
	}
	return parsing.parseInstructionOnExpression(beginOffset, leftHandSide)
}

func (parsing *Parsing) parseInstructionOnExpression(
	beginOffset input.Offset, leftHandSide tree.Expression) (tree.Expression, error) {

	switch operator := token.OperatorValue(parsing.token()); {
	case operator.IsAssign():
		parsing.advance() // TODO: ?
		return parsing.parseAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
	case operator == token.DecrementOperator:
		return parsing.parsePostfixExpressionOnNode(beginOffset, leftHandSide)
	}
	parsing.advance()
	return &tree.ExpressionStatement{
		Expression: leftHandSide,
	}, nil
}

func (parsing *Parsing) parsePostfixExpressionOnNode(
	beginOffset input.Offset, leftHandSide tree.Expression) (tree.Expression, error) {

	operation := token.OperatorValue(parsing.token())
	parsing.advance()
	parsing.skipEndOfStatement()
	return &tree.PostfixExpression{
		Operand:  leftHandSide,
		Operator: operation,
		Region:   parsing.createRegion(beginOffset),
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
		Target:   declaration,
		Value:    value,
		Operator: token.AssignOperator,
		Region:   parsing.createRegion(beginOffset),
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
		Name:     fieldName,
		TypeName: typeName,
		Region:   parsing.createRegion(beginOffset),
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
		Target:   declaration,
		Value:    value,
		Operator: token.AssignOperator,
		Region:   parsing.createRegion(beginOffset),
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
		Value:  baseTypeOrAccessedField.Value(),
		Region: parsing.createRegion(beginOffset),
	})
	if err != nil {
		return nil, err
	}
	return parsing.parseOperationOrAssign(operation)
}

func (parsing *Parsing) parseFieldDeclarationFromBaseTypeName(
	beginOffset input.Offset, baseTypeName token.Token) (tree.Node, error) {

	typeName, err := parsing.parseTypeNameFromBaseIdentifier(beginOffset, baseTypeName)
	if err != nil {
		return nil, err
	}
	return parsing.parseFieldDeclarationWithTypeName(beginOffset, typeName)
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
func (parsing *Parsing) parseStatementSequence() (statements []tree.Statement) {
	for {
		expectedIndent := parsing.block.Indent
		current := parsing.token()
		if token.IsEndOfFileToken(current) {
			break
		}
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
func (parsing *Parsing) parseStatementBlock() *tree.BlockStatement {
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
		Children: statements,
		Region:   parsing.createRegion(beginOffset),
	}, nil
}
