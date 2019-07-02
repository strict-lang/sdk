package parser

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"github.com/BenjaminNitschke/Strict/pkg/token"
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

// UnexpectedTokenError indicates that the parser expected a certain kind of token, but
// got a different one. It captures the token and has an optional 'expected' field, which
// stores the name of the kind of token that was expected.
type UnexpectedTokenError struct {
	token    token.Token
	expected string
}

func (err *UnexpectedTokenError) Error() string {
	if err.expected != "" {
		return fmt.Sprintf("expected %s but got %s", err.expected, err.token)
	}
	return fmt.Sprintf("unexpected token: %s", err.token)
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