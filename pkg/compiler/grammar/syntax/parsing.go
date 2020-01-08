package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/code"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"log"
)

var notParsingMethod = parsedMethod{name: `!none`}

// Parsing represents the process of grammar a stream of tokens and turning them
// into an abstract grammar tree. This class is not reusable and can only produce
// one translation unit. It does some scope management but does not do to many
// checks that could be considered semantic.
type Parsing struct {
	tokenReader token.Stream
	rootScope   *code.Scope
	recorder    *diagnostic.Bag
	block       *Block
	unitName    string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the grammar looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
	// name of the method which is currently parsed. It is required by the optional
	// test statement within a method. The name is set to an empty string after a
	// method has been parsed.
	currentMethod        parsedMethod
	isAtBeginOfStatement bool
	structureStack       *structureStack
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the grammar to scanning code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Scope  *code.Scope
	Parent *Block
}

func (parsing *Parsing) parseImportStatementList() (imports []*tree.ImportStatement) {
	for token.HasKeywordValue(parsing.token(), token.ImportKeyword) {
		imports = append(imports, parsing.parseImportStatement())
	}
	return imports
}

func (parsing *Parsing) parseClassDeclaration() *tree.ClassDeclaration {
	parsing.beginStructure(tree.ClassDeclarationNodeKind)
	nodes := parsing.parseTopLevelNodes()
	return &tree.ClassDeclaration{
		Name:       parsing.unitName,
		Parameters: []*tree.ClassParameter{},
		SuperTypes: []tree.TypeName{},
		Children:   nodes,
		Region:     parsing.completeStructure(tree.ClassDeclarationNodeKind),
	}
}

// Parse parses a TranslationUnit and returns an error on failure.
func (parsing *Parsing) Parse() (result *tree.TranslationUnit, err error) {
	defer func() {
		if failure := recover(); failure != nil {
			err = extractErrorFromPanic(failure)
		}
	}()
	return parsing.parseTranslationUnit(), nil
}

// parseTranslationUnit invokes the grammar on the translation unit.
// This method can only be called once on the Parsing instance.
func (parsing *Parsing) parseTranslationUnit() *tree.TranslationUnit {
	parsing.beginStructure(tree.TranslationUnitNodeKind)
	imports := parsing.parseImportStatementList()
	class := parsing.parseClassDeclaration()
	return &tree.TranslationUnit{
		Name:    parsing.unitName,
		Imports: imports,
		Class:   class,
		Region:  parsing.completeStructure(tree.TranslationUnitNodeKind),
	}
}

// openBlock opens a new block of code, updates the grammar block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parsing *Parsing) openBlock(indent token.Indent) {
	var blockScope *code.Scope
	if parsing.block == nil {
		blockScope = parsing.rootScope.NewChild()
	} else {
		blockScope = parsing.block.Scope.NewChild()
	}
	block := &Block{
		Indent: indent,
		Scope:  blockScope,
		Parent: parsing.block,
	}
	parsing.block = block
}

func (parsing *Parsing) token() token.Token {
	return parsing.tokenReader.Last()
}

func (parsing *Parsing) pullToken() token.Token {
	last := parsing.tokenReader.Last()
	parsing.tokenReader.Pull()
	return last
}

func (parsing *Parsing) advance() {
	parsing.tokenReader.Pull()
	parsing.isAtBeginOfStatement = false
}

func (parsing *Parsing) peek() token.Token {
	return parsing.tokenReader.Peek()
}

func (parsing *Parsing) closeBlock() {
	parsing.block = parsing.block.Parent
}

func (parsing *Parsing) offset() input.Offset {
	return parsing.token().Position().Begin()
}

func (parsing *Parsing) beginStructure(kind tree.NodeKind) {
	parsing.pushStructure(structureStackElement{
		beginOffset: parsing.offset(),
		nodeKind:    kind,
	})
}

func (parsing *Parsing) updateTopStructureKind(kind tree.NodeKind) {
	// Stack elements are values thus we can not change the top's field.
	// Instead the top element is exchanged.
	top, err := parsing.structureStack.pop()
	if err == nil {
		top.nodeKind = kind
		parsing.pushStructure(top)
	}
}

func (parsing *Parsing) completeStructure(expectedKind tree.NodeKind) input.Region {
	structure, err := parsing.structureStack.pop()
	if err != nil {
		parsing.throwError(newEmptyStructureStackError(expectedKind))
		return input.Region{}
	}
	if expectedKind != tree.WildcardNodeKind && structure.nodeKind != expectedKind {
		log.Printf("Expected to complete %s but completed %s", expectedKind, structure.nodeKind)
	}
	begin := structure.beginOffset
	return input.CreateRegion(begin, parsing.offset())
}

func newEmptyStructureStackError(expected tree.NodeKind) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error:         &diagnostic.InvalidStatementError{Kind: expected},
		CommonReasons: []string{"Internal bug in the compiler"},
	}
}

func (parsing *Parsing) createRegionOfCurrentStructure() input.Region {
	begin := parsing.peekStructure().beginOffset
	return input.CreateRegion(begin, parsing.offset())
}

func (parsing *Parsing) pushStructure(structure structureStackElement) {
	parsing.structureStack.push(structure)
}

func (parsing *Parsing) peekStructure() structureStackElement {
	return parsing.structureStack.peek()
}

func (parsing *Parsing) isParsingMethod() bool {
	return parsing.currentMethod != notParsingMethod
}

func (parsing *Parsing) parseTopLevelNodes() (nodes []tree.Node) {
	parsing.beginStructure(tree.StatementBlockNodeKind)
	block := parsing.parseStatementBlock()
	defer func() {
		if failure := recover(); failure != nil {
			err := extractErrorFromPanic(failure)
			invalid := parsing.completeInvalidStructure(err)
			nodes = []tree.Node{invalid}
		}
	}()
	parsing.completeStructure(tree.StatementBlockNodeKind)
	return convertStatementSliceToNodeSlice(block.Children)
}

func extractErrorFromPanic(value interface{}) error {
	if err, isError := value.(error); isError {
		return err
	}
	return fmt.Errorf("%s", value)
}

func convertStatementSliceToNodeSlice(statements []tree.Statement) (nodes []tree.Node) {
	for _, statement := range statements {
		nodes = append(nodes, statement)
	}
	return nodes
}

// skipOperator skips the next keyword if it the passed operator, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipOperator(operator token.Operator) {
	if err := parsing.expectOperator(operator); err != nil {
		parsing.throwError(err)
	}
	parsing.advance()
}

// skipKeyword skips the next keyword if it the passed keyword, otherwise
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) skipKeyword(keyword token.Keyword) {
	if err := parsing.expectKeyword(keyword); err != nil {
		parsing.throwError(err)
	}
	parsing.advance()
}

// expectOperator peeks the next token and expects it to be the passed operator,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectOperator(expected token.Operator) *diagnostic.RichError {
	if token.OperatorValue(parsing.token()) != expected {
		return newInvalidOperatorError(parsing.token(), expected)
	}
	return nil
}

// expectKeyword peeks the next token and expects it to be the passed keyword,
// otherwise an UnexpectedTokenError is returned.
func (parsing *Parsing) expectKeyword(expected token.Keyword) *diagnostic.RichError {
	if token.KeywordValue(parsing.token()) != expected {
		return newInvalidKeywordError(parsing.token(), expected)
	}
	return nil
}

// expectAnyIdentifier expects the next token to be an identifier,
// without regards to its value and returns an error if it fails.
func (parsing *Parsing) expectAnyIdentifier() *tree.Identifier {
	parsing.beginStructure(tree.IdentifierNodeKind)
	current := parsing.token()
	if !token.IsIdentifierToken(current) {
		parsing.throwError(newNoIdentifierError(current))
	}
	return &tree.Identifier{
		Value:  current.Value(),
		Region: parsing.completeStructure(tree.IdentifierNodeKind),
	}
}

func (parsing *Parsing) isLookingAtKeyword(keyword token.Keyword) bool {
	return token.HasKeywordValue(parsing.peek(), keyword)
}

func (parsing *Parsing) isLookingAtOperator(operator token.Operator) bool {
	return token.HasOperatorValue(parsing.peek(), operator)
}

func (parsing *Parsing) completeInvalidStructure(err error) tree.Statement {
	region := parsing.completeStructure(tree.WildcardNodeKind)
	parsing.reportError(newInvalidStructureError(), region)
	return &tree.InvalidStatement{
		Region: region,
	}
}

func newInvalidStructureError() *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.InvalidStatementError{Kind: tree.UnknownNodeKind},
	}
}

// skipEndOfStatement skips the next token if it is an EndOfStatement token.
func (parsing *Parsing) skipEndOfStatement() {
	// Do not report the missing end of statement.
	parsing.advance()
	parsing.isAtBeginOfStatement = true
}

// reportError reports an error to the diagnostics bag, starting at the
// passed position and ending at the parsers current position.
func (parsing *Parsing) reportError(error *diagnostic.RichError, region input.Region) {
	parsing.recorder.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SyntacticalAnalysis,
		UnitName: parsing.unitName,
		Position: region,
		Error:    error,
	})
}

type parsingError struct {
	Structure tree.NodeKind
	Position  input.Region
	Cause     error
}

func (err *parsingError) String() string {
	return fmt.Sprintf("failed parsing %s: %s",
		err.Structure.Name(),
		err.Cause.Error())
}

func (parsing *Parsing) throwError(cause *diagnostic.RichError) {
	structure, err := parsing.structureStack.pop()
	if err != nil {
		structure = structureStackElement{nodeKind: tree.UnknownNodeKind}
		log.Print("Could not pop structure stack")
	}
	region := input.CreateRegion(structure.beginOffset, parsing.offset())
	parsing.reportError(cause, region)
	panic(&parsingError{
		Structure: structure.nodeKind,
		Position:  region,
		Cause:     fmt.Errorf("could not parse %s", structure.nodeKind),
	})
}
func newNoIdentifierError(token token.Token) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "Identifier",
			Received: token.Value(),
		},
		CommonReasons: []string{
			"Declarations are not written properly",
		},
	}
}

func newInvalidKeywordError(token token.Token, expected token.Keyword) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: expected.String(),
			Received: token.Value(),
		},
		CommonReasons: aggregateReasonsOfInvalidKeywordError(token, expected),
	}
}

const forgottenDoKeywordReason = "The 'do' keyword after control structures is forgotten"

func aggregateReasonsOfInvalidKeywordError(
	received token.Token,
	expected token.Keyword) (reasons []string) {

	if expected == token.DoKeyword && token.IsEndOfStatementToken(received) {
		reasons = append(reasons, forgottenDoKeywordReason)
	}
	return reasons
}

func newInvalidOperatorError(token token.Token, expected token.Operator) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: expected.String(),
			Received: token.Value(),
		},
		CommonReasons: aggregateReasonsOfInvalidOperatorError(token),
	}
}

const unfinishedOperationReason = "An operation has not been completed"
const invalidOperator = "An invalid operator is applied to an operation"

func aggregateReasonsOfInvalidOperatorError(received token.Token) (reasons []string) {

	if token.IsIdentifierToken(received) {
		reasons = append(reasons, unfinishedOperationReason)
	}
	reasons = append(reasons, invalidOperator)
	return reasons
}
