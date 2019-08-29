package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/scope"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

const notParsingMethod = ""

// Parsing represents the process of parsing a stream of tokens and turning them
// into an abstract syntax tree. This class is not reusable and can only produce
// one translation unit. It does some scope management but does not do to many
// checks that could be considered semantic.
type Parsing struct {
	tokenReader token.Reader
	rootScope   *scope.Scope
	recorder    *diagnostic.Bag
	block       *Block
	unitName    string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the parsing looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
	// name of the method which is currently parsed. It is required by the optional
	// test statement within a method. The name is set to an empty string after a
	// method has been parsed.
	currentMethodName string
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the parsing to scanning code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Scope  *scope.Scope
	Parent *Block
}

// openBlock opens a new block of code, updates the parsing block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parsing *Parsing) openBlock(indent token.Indent) {
	var blockScope *scope.Scope
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

func (parsing *Parsing) reportError(err error, position ast.Position) {
	parsing.recorder.Record(diagnostic.RecordedEntry{
		Kind:    &diagnostic.Error,
		Stage:   &diagnostic.SyntacticalAnalysis,
		Message: err.Error(),
		Position: position,
	})
}

func (parsing *Parsing) token() token.Token {
	return parsing.tokenReader.Last()
}

func (parsing *Parsing) advance() {
	parsing.tokenReader.Pull()
}

func (parsing *Parsing) peek() token.Token {
	return parsing.tokenReader.Peek()
}

func (parsing *Parsing) closeBlock() {
	parsing.block = parsing.block.Parent
}

func (parsing *Parsing) offset() source.Offset {
	return parsing.token().Position().Begin()
}

type offsetPosition struct {
	begin source.Offset
	end   source.Offset
}

func (position offsetPosition) Begin() source.Offset {
	return position.begin
}

func (position offsetPosition) End() source.Offset {
	return position.end
}

func (parsing *Parsing) createPosition(beginOffset source.Offset) ast.Position {
	return &offsetPosition{begin: beginOffset, end: parsing.offset()}
}

func (parsing *Parsing) createTokenPosition() ast.Position {
	return parsing.token().Position()
}

func (parsing *Parsing) parseTopLevelNodes() []ast.Node {
	beginOffset := parsing.offset()
	block, err := parsing.ParseStatementBlock()
	if err != nil {
		return []ast.Node{parsing.createInvalidStatement(beginOffset, err)}
	}
	return block.Children
}

func (parsing *Parsing) ParseTranslationUnit() (*ast.TranslationUnit, error) {
	topLevelNodes := parsing.parseTopLevelNodes()
	return ast.NewTranslationUnit(parsing.unitName, parsing.rootScope, topLevelNodes), nil
}

func (parsing *Parsing) isParsingMethod() bool {
	return parsing.currentMethodName != notParsingMethod
}
