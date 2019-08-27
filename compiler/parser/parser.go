package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/scope"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

const notParsingMethod = ""

// Parser parses an AST from a stream of tokens.
type Parser struct {
	tokenReader token.Reader
	rootScope   *scope.Scope
	recorder    *diagnostic.Recorder
	block       *Block
	unitName    string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the parser looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
	// name of the method which is currently parsed. It is required by the optional
	// test statement within a method. The name is set to an empty string after a
	// method has been parsed.
	currentMethodName string
}

// Block represents a nested sequence of statements that has a set indentation level.
// It helps the parser to scan code blocks and know where a block ends.
type Block struct {
	Indent token.Indent
	Scope  *scope.Scope
	Parent *Block
}

// openBlock opens a new block of code, updates the parser block pointer and
// creates a new scope for that block that is a child-scope of the parsers
// last block. Only statements with the blocks indent may go into the block.
func (parser *Parser) openBlock(indent token.Indent) {
	var blockScope *scope.Scope
	if parser.block == nil {
		blockScope = parser.rootScope.NewChild()
	} else {
		blockScope = parser.block.Scope.NewChild()
	}
	block := &Block{
		Indent: indent,
		Scope:  blockScope,
		Parent: parser.block,
	}
	parser.block = block
}

func (parser *Parser) reportError(err error) {
	parser.recorder.Record(diagnostic.RecordedEntry{
		Kind:    &diagnostic.Error,
		Stage:   &diagnostic.SyntacticalAnalysis,
		Message: err.Error(),
	})
}

func (parser *Parser) token() token.Token {
	return parser.tokenReader.Last()
}

func (parser *Parser) advance() {
	parser.tokenReader.Pull()
}

func (parser *Parser) peek() token.Token {
	return parser.tokenReader.Peek()
}

func (parser *Parser) closeBlock() {
	parser.block = parser.block.Parent
}

func (parser *Parser) offset() source.Offset {
	return parser.token().Position().Begin
}

type offsetPosition struct {
	begin source.Offset
	end source.Offset
}

func (position offsetPosition) Begin() source.Offset {
	return position.begin
}

func (position offsetPosition) End() source.Offset {
	return position.end
}

func (parser *Parser) createPosition(beginOffset source.Offset) ast.Position {
	return offsetPosition{begin: beginOffset, end: parser.offset()}
}

func (parser *Parser) createTokenPosition() ast.Position { return nil }

func (parser *Parser) parseTopLevelNodes() []ast.Node {
	beginOffset := parser.offset()
	block, err := parser.ParseStatementBlock()
	if err != nil {
		return []ast.Node{parser.createInvalidStatement(beginOffset, err)}
	}
	return block.Children
}

func (parser *Parser) ParseTranslationUnit() (*ast.TranslationUnit, error) {
	topLevelNodes := parser.parseTopLevelNodes()
	return ast.NewTranslationUnit(parser.unitName, parser.rootScope, topLevelNodes), nil
}

func (parser *Parser) isParsingMethod() bool {
	return parser.currentMethodName != notParsingMethod
}
