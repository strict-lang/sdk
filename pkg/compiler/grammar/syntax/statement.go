package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() *tree.ConditionalStatement {
	parsing.beginStructure(tree.ConditionalStatementNodeKind)
	parsing.skipKeyword(token.IfKeyword)
	condition := parsing.parseConditionalExpression()
	parsing.skipKeyword(token.DoKeyword)
	parsing.skipEndOfStatement()
	consequence := parsing.parseStatementBlock()
	return parsing.parseElseClauseIfPresent(condition, consequence)
}

func (parsing *Parsing) parseElseClauseIfPresent(
	condition tree.Expression,
	consequence tree.Statement) *tree.ConditionalStatement {

	if token.HasKeywordValue(parsing.token(), token.ElseKeyword) {
		return parsing.parseConditionalStatementWithAlternative(
			condition, consequence)
	}
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Region:      parsing.completeStructure(tree.ConditionalStatementNodeKind),
	}
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	condition tree.Node,
	consequence tree.Node) *tree.ConditionalStatement {

	parsing.advance()
	alternative := parsing.parseElseIfOrBlock()
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
		Region:      parsing.completeStructure(tree.ConditionalStatementNodeKind),
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
	parsing.beginStructure(tree.ForEachLoopStatementNodeKind)
	parsing.skipKeyword(token.ForKeyword)
	if parsing.isLookingAtRangedLoop() {
		return parsing.completeFromToStatement()
	}
	return parsing.completeForEachStatement()
}

func (parsing *Parsing) isLookingAtRangedLoop() bool {
	return token.IsIdentifierToken(parsing.token()) &&
		token.HasKeywordValue(parsing.peek(), token.FromKeyword)
}

// completeForEachStatement is called by the ParseForStatement method
// after it checked for a foreach statement. At this point the last token
// is an identifier that is the name of the foreach loops element field.
// This method completes the loops grammar.
func (parsing *Parsing) completeForEachStatement() *tree.ForEachLoopStatement {
	parsing.updateTopStructureKind(tree.ForEachLoopStatementNodeKind)
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
		Region:   parsing.completeStructure(tree.ForEachLoopStatementNodeKind),
	}
}

// completeFromToStatement is called by the ParseForStatement method
// after it peeked the 'from' keyword. At this point, the last token
// is an identifier that is the name of the loops counter field. This
// method completes the loops grammar.
func (parsing *Parsing) completeFromToStatement() *tree.RangedLoopStatement {
	parsing.updateTopStructureKind(tree.RangedLoopStatementNodeKind)
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
		Region: parsing.completeStructure(tree.RangedLoopStatementNodeKind),
	}
}

func (parsing *Parsing) parseYieldStatement() *tree.YieldStatement {
	parsing.beginStructure(tree.YieldStatementNodeKind)
	parsing.skipKeyword(token.YieldKeyword)
	rightHandSide := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.YieldStatement{
		Value:  rightHandSide,
		Region: parsing.completeStructure(tree.YieldStatementNodeKind),
	}
}

func (parsing *Parsing) parseReturnStatement() *tree.ReturnStatement {
	parsing.beginStructure(tree.ReturnStatementNodeKind)
	parsing.skipKeyword(token.ReturnKeyword)
	defer parsing.skipEndOfStatement()
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &tree.ReturnStatement{
			Region: parsing.completeStructure(tree.ReturnStatementNodeKind),
		}
	}
	rightHandSide := parsing.parseExpression()
	return &tree.ReturnStatement{
		Value:  rightHandSide,
		Region: parsing.completeStructure(tree.ReturnStatementNodeKind),
	}
}

func (parsing *Parsing) parseNestedMethodDeclaration() *tree.MethodDeclaration {
	return parsing.parseMethodDeclaration()
}

func (parsing *Parsing) parseImportStatement() *tree.ImportStatement {
	parsing.beginStructure(tree.ImportStatementNodeKind)
	parsing.skipKeyword(token.ImportKeyword)
	switch next := parsing.token(); {
	case token.IsStringLiteralToken(next):
		return parsing.completeFileImport()
	case token.IsIdentifierToken(next):
		return parsing.completeIdentifierChainImport()
	default:
		parsing.throwError(&UnexpectedTokenError{
			Expected: "File or Class Path",
			Token:    next,
		})
		return nil
	}
}

func (parsing *Parsing) completeIdentifierChainImport() *tree.ImportStatement {
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
		Region: parsing.completeStructure(tree.ImportStatementNodeKind),
	}
}

func (parsing *Parsing) completeFileImport() *tree.ImportStatement {
	target := &tree.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &tree.ImportStatement{
			Target: target,
			Region: parsing.completeStructure(tree.ImportStatementNodeKind),
		}
	}
	parsing.advance()
	return parsing.completeFileImportWithAlias(target)
}

func (parsing *Parsing) completeFileImportWithAlias(target tree.ImportTarget) *tree.ImportStatement {
	alias := parsing.parseImportAlias()
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target: target,
		Alias:  alias,
		Region: parsing.completeStructure(tree.ImportStatementNodeKind),
	}
}

func (parsing *Parsing) parseImportAlias() *tree.Identifier {
	return parsing.parseIdentifier()
}

func (parsing *Parsing) parseOptionalAssignValue() tree.Expression {
	parsing.skipOperator(token.AssignOperator)
	return parsing.parseExpression()
}

type keywordStatementParser func(*Parsing) tree.Node

var keywordStatementParserTable map[token.Keyword]keywordStatementParser

func init() {
	keywordStatementParserTable = map[token.Keyword]keywordStatementParser{
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

func (parsing *Parsing) maybeParseConstructorDeclaration() (keywordStatementParser, bool) {
	if !parsing.shouldParseConstructorDeclaration() {
		return nil, false
	}
	return func(*Parsing) tree.Node {
		return parsing.parseConstructorDeclaration()
	}, true
}

func (parsing *Parsing) parseConstructorDeclaration() tree.Node {
	parsing.beginStructure(tree.ConstructorDeclarationNodeKind)
	parsing.skipKeyword(token.CreateKeyword)
	parameters := parsing.parseParameterList()
	parsing.skipEndOfStatement()
	body := parsing.parseStatementBlock()
	return &tree.ConstructorDeclaration{
		Parameters: parameters,
		Child:      body,
		Region:     parsing.completeStructure(tree.ConstructorDeclarationNodeKind),
	}
}

// parseKeywordStatement parses a statement that starts with a keyword.
func (parsing *Parsing) parseKeywordStatement(keyword token.Keyword) tree.Node {
	function, ok := parsing.findKeywordStatementParser(keyword)
	if ok {
		return function(parsing)
	}
	parsing.throwError(&UnexpectedTokenError{
		Token:    parsing.token(),
		Expected: "keyword that starts a statement",
	})
	return nil
}

func (parsing *Parsing) completeAssignStatement(
	operator token.Operator, leftHandSide tree.Node) tree.Node {

	parsing.updateTopStructureKind(tree.AssignStatementNodeKind)
	rightHandSide := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.AssignStatement{
		Target:   leftHandSide,
		Value:    rightHandSide,
		Operator: operator,
		Region:   parsing.completeStructure(tree.AssignStatementNodeKind),
	}
}

func (parsing *Parsing) parseTestStatement() tree.Node {
	parsing.beginStructure(tree.TestStatementNodeKind)
	parsing.skipKeyword(token.TestKeyword)
	parsing.skipEndOfStatement()
	statements := parsing.parseStatementBlock()
	return &tree.TestStatement{
		Region:     parsing.completeStructure(tree.TestStatementNodeKind),
		MethodName: parsing.currentMethod.name,
		Child:      statements,
	}
}

func (parsing *Parsing) parseAssertStatement() tree.Node {
	parsing.beginStructure(tree.AssertStatementNodeKind)
	parsing.skipKeyword(token.AssertKeyword)
	expression := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.AssertStatement{
		Region:     parsing.completeStructure(tree.AssertStatementNodeKind),
		Expression: expression,
	}
}

// parseInstructionStatement parses a statement that is not a structured-control flow
// statement. Instructions mostly operate on values and assign fields.
func (parsing *Parsing) parseInstructionStatement() tree.Node {
	parsing.beginStructure(tree.UnknownNodeKind)
	leftHandSide := parsing.parseExpression()
	return parsing.completeInstructionOnNode(leftHandSide)
}

func (parsing *Parsing) completeInstructionOnNode(
	leftHandSide tree.Expression) tree.Expression {

	switch operator := token.OperatorValue(parsing.token()); {
	case operator.IsAssign():
		parsing.advance()
		return parsing.completeAssignStatement(operator, leftHandSide)
	case operator == token.IncrementOperator:
	case operator == token.DecrementOperator:
		return parsing.completePostfixExpressionOnNode(leftHandSide)
	}
	parsing.advance()
	parsing.completeStructure(tree.WildcardNodeKind)
	return &tree.ExpressionStatement{
		Expression: leftHandSide,
	}
}

func (parsing *Parsing) completePostfixExpressionOnNode(
	leftHandSide tree.Expression) tree.Expression {

	operation := token.OperatorValue(parsing.token())
	parsing.advance()
	parsing.skipEndOfStatement()
	return &tree.PostfixExpression{
		Operand:  leftHandSide,
		Operator: operation,
		Region:   parsing.completeStructure(tree.PostfixExpressionNodeKind),
	}
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
	return parsing.couldBeLookingAtTypeName()
}

func (parsing *Parsing) parseFieldDeclarationOrDefinition() tree.Node {
	parsing.beginStructure(tree.FieldDeclarationNodeKind)
	declaration := parsing.parseFieldDeclaration()
	if !token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		parsing.skipEndOfStatement()
		return declaration
	}
	parsing.advance()
	return parsing.completeParsingFieldDefinition(declaration)
}
func (parsing *Parsing) completeParsingFieldDefinition(
	declaration tree.Node) tree.Node {

	parsing.updateTopStructureKind(tree.AssignStatementNodeKind)
	value := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.AssignStatement{
		Target:   declaration,
		Value:    value,
		Operator: token.AssignOperator,
		Region:   parsing.completeStructure(tree.AssignStatementNodeKind),
	}
}

func (parsing *Parsing) parseFieldDeclaration() tree.Node {
	parsing.beginStructure(tree.FieldDeclarationNodeKind)
	typeName := parsing.parseTypeName()
	return parsing.completeFieldDeclarationWithTypeName(typeName)
}

func (parsing *Parsing) completeFieldDeclarationWithTypeName(
	typeName tree.TypeName) tree.Node {

	parsing.updateTopStructureKind(tree.FieldDeclarationNodeKind)
	fieldName := parsing.parseIdentifier()
	declaration := &tree.FieldDeclaration{
		Name:     fieldName,
		TypeName: typeName,
		Region:   parsing.completeStructure(tree.FieldDeclarationNodeKind),
	}
	if token.HasOperatorValue(parsing.token(), token.AssignOperator) {
		return parsing.completeFieldDefinition(declaration)
	}
	return declaration
}

func (parsing *Parsing) completeFieldDefinition(declaration *tree.FieldDeclaration) tree.Node {
	parsing.skipOperator(token.AssignOperator)
	value := parsing.parseExpression()
	return &tree.AssignStatement{
		Target:   declaration,
		Value:    value,
		Operator: token.AssignOperator,
		Region:   parsing.completeStructure(tree.FieldDeclarationNodeKind),
	}
}

// ParseStatement parses the next statement from the stream of tokens. Statements include
// conditionals or loops, therefor this function may end up scanning multiple statements
// and call itself.
func (parsing *Parsing) parseStatement() tree.Node {
	parsing.beginStructure(tree.UnknownNodeKind)
	defer parsing.completeStructure(tree.WildcardNodeKind) // Could have been modified

	switch current := parsing.token(); {
	case parsing.isKeywordStatementToken(current):
		return parsing.parseKeywordStatement(token.KeywordValue(current))
	case parsing.isKeywordExpressionToken(current):
		fallthrough
	case token.IsIdentifierToken(current):
		if parsing.shouldParseFieldDeclaration() {
			return parsing.parseFieldDeclarationOrListAccess()
		}
		fallthrough
	case token.IsOperatorToken(current):
		fallthrough
	case token.IsLiteralToken(current):
		return parsing.parseInstructionStatement()
	default:
		parsing.throwError(fmt.Errorf("expected begin of statement or expression but got: %s", current))
		return nil
	}
}

// parseFieldDeclarationOrListAccess either parses a FieldDeclaration or a ListSelect.
// Both constructs are similar when only being able to peek one token.
func (parsing *Parsing) parseFieldDeclarationOrListAccess() tree.Node {
	parsing.beginStructure(tree.UnknownNodeKind)
	baseTypeOrAccessedField := parsing.pullToken()
	if parsing.isLookingAtListAccess() {
		return parsing.completeFieldDeclarationFromBaseTypeName(baseTypeOrAccessedField)
	}
	return parsing.completeListAccess(baseTypeOrAccessedField)
}

func createRegionForToken(target token.Token) input.Region {
	position := target.Position()
	return input.CreateRegion(position.Begin(), position.End())
}

func (parsing *Parsing) completeListAccess(target token.Token) tree.Statement {
	field := &tree.Identifier{
		Region: createRegionForToken(target),
		Value:  target.Value(),
	}
	operation := parsing.parseOperationsOnOperand(field)
	// This method has to be used because of selector chaining and calls.
	return parsing.parseOperationOrAssign(operation)
}

func (parsing *Parsing) isLookingAtListAccess() bool {
	return token.HasOperatorValue(parsing.token(), token.LeftBracketOperator) &&
		token.HasOperatorValue(parsing.peek(), token.RightBracketOperator)
}

func (parsing *Parsing) completeFieldDeclarationFromBaseTypeName(
	baseTypeName token.Token) tree.Node {

	typeName := parsing.parseTypeNameFromBaseIdentifier(baseTypeName.Value())
	return parsing.completeFieldDeclarationWithTypeName(typeName)
}

// ParseStatementSequence parses a sequence of statements. The sequence
// is ended when the first token in a line has an indent other than the
// value in the current blocks indent field.
// TODO: Maybe recover from error
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
			parsing.throwError(fmt.Errorf("indent level of %d", expectedIndent))
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
	parsing.beginStructure(tree.BlockStatementNodeKind)
	indent := parsing.token().Indent()
	if indent < parsing.block.Indent {
		parsing.throwError(&InvalidIndentationError{
			Token:    parsing.token(),
			Expected: "indent bigger than 0",
		})
	}
	parsing.openBlock(indent)
	statements := parsing.parseStatementSequence()
	parsing.closeBlock()
	return &tree.BlockStatement{
		Children: statements,
		Region:   parsing.completeStructure(tree.BlockStatementNodeKind),
	}
}
