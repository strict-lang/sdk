package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/scope"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

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
	offset := parser.token().Position().Begin
	parser.recorder.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SyntacticalAnalysis,
		Source:   parser.token().Value(),
		Message:  err.Error(),
		Offset: offset,
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

func (parser *Parser) parseTopLevelNodes() []ast.Node {
	block, err := parser.ParseStatementBlock()
	if err != nil {
		return []ast.Node{parser.createInvalidStatement(err)}
	}
	return block.Children
}

func (parser *Parser) ParseTranslationUnit() (*ast.TranslationUnit, error) {
	topLevelNodes := parser.parseTopLevelNodes()
	return ast.NewTranslationUnit(parser.unitName, parser.rootScope, topLevelNodes), nil
}