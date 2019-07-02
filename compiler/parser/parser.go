package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scope"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

// Parser parses an AST from a stream of tokens.
type Parser struct {
	tokens   token.Reader
	unit     *ast.TranslationUnit
	recorder *diagnostic.Recorder
	block		 *Block
}

// NewParser creates a parser instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func NewParser(unit *ast.TranslationUnit, tokens token.Reader, recorder *diagnostic.Recorder) *Parser {
	return &Parser{
		unit:     unit,
		tokens:   tokens,
		recorder: recorder,
	}
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
		blockScope = parser.unit.Scope().NewChild()
	} else {
		blockScope = parser.block.Scope.NewChild()
	}
	block := &Block{
		Indent: indent,
		Scope: blockScope,
		Parent: parser.block,
	}
	parser.block = block
}

func (parser *Parser) closeBlock() {
	parser.block = parser.block.Parent
}