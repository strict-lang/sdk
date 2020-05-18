package syntax

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var notParsingMethod = parsedMethod{name: `!none`}

type Parsing struct {
	tokenReader     token.StreamWithLineMap
	recorder        *diagnostic.Bag
	block           *Block
	unitName        string
	expressionDepth int
	currentMethod   parsedMethod
	statementBegin  bool
	structureStack  *structureStack
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the grammar to scanning code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Parent *Block
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

func (parsing *Parsing) parseTranslationUnit() *tree.TranslationUnit {
	parsing.beginStructure(tree.TranslationUnitNodeKind)
	imports := parsing.parseImportStatementList()
	class := parsing.parseClassDeclaration()
	return &tree.TranslationUnit{
		Name:    parsing.unitName,
		Imports: imports,
		Class:   class,
		LineMap: parsing.tokenReader.NewLineMap(),
		Region:  parsing.completeStructure(tree.TranslationUnitNodeKind),
	}
}


// openBlock opens a new block of code, updates the grammar block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parsing *Parsing) openBlock(indent token.Indent) {
	block := &Block{
		Indent: indent,
		Parent: parsing.block,
	}
	parsing.block = block
}

func (parsing *Parsing) closeBlock() {
	parsing.block = parsing.block.Parent
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
	parsing.statementBegin = false
}

func (parsing *Parsing) peek() token.Token {
	return parsing.tokenReader.Peek()
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

func convertFileNameToClassName(fileName string) string {
	nameWithoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	lastDot := strings.LastIndex(nameWithoutExtension, string(os.PathSeparator))
	if lastDot != -1 {
		return nameWithoutExtension[lastDot + 1:]

	}
	return nameWithoutExtension
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

func (parsing *Parsing) parseTopLevelDeclarations() (nodes []tree.Statement) {
	for {
		current := parsing.token()
		if token.IsEndOfFileToken(current) {
			break
		}
		if token.IsEndOfStatementToken(current) {
			parsing.advance()
			continue
		}
		nodes = append(nodes, parsing.parseTopLevelDeclaration())
	}
	return
}

func (parsing *Parsing) parseTopLevelDeclaration() tree.Statement {
	current := parsing.token()
	if token.IsKeywordToken(current) {
		return parsing.parseKeywordStatement(token.KeywordValue(current))
	} else {
		parsing.throwError(newUnexpectedTokenError(current))
		return &tree.InvalidStatement{}
	}
}

func (parsing *Parsing) parseTopLevelNodes() (nodes []tree.Node) {
	block := parsing.parseTopLevelDeclarations()
	defer func() {
		if failure := recover(); failure != nil {
			err := extractErrorFromPanic(failure)
			invalid := parsing.completeInvalidStructure(err)
			nodes = []tree.Node{invalid}
		}
	}()
	return convertStatementSliceToNodeSlice(block)
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
