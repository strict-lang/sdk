package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/scope"
	"github.com/BenjaminNitschke/Strict/compiler/source/linemap"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

// Parser parses an AST from a stream of tokens.
type Parser struct {
	tokens    token.Reader
	rootScope *scope.Scope
	recorder  *diagnostic.Recorder
	linemap   *linemap.Linemap
	block     *Block
	unitName  string
	// expressionDepth is the amount of parentheses encountered at the
	// current time. It is incremented every time the parser looks at a
	// LeftParenOperator and decremented when it looks at a RightParenOperator.
	expressionDepth int
}

// NewParser creates a parser instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func NewParser(unitName string, tokens token.Reader, recorder *diagnostic.Recorder) *Parser {
	parser := &Parser{
		rootScope: scope.NewRoot(),
		tokens:    tokens,
		recorder:  recorder,
		unitName:  unitName,
	}
	parser.openBlock(token.NoIndent)
	return parser
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
	parser.recorder.Record(diagnostic.Entry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SyntacticalAnalysis,
		Source:   parser.tokens.Peek().Value(),
		Message:  err.Error(),
		Position: diagnostic.Position{
			// TODO(merlinosayimwen): Use linemap to get line information of
			// 	the token and create a diagnostic.Position from it.
		},
	})
}

func (parser *Parser) closeBlock() {
	parser.block = parser.block.Parent
}

func (parser *Parser) parseTopLevelNodes() ([]ast.Node, error) {
	return []ast.Node{}, nil
}

func (parser *Parser) ParseTranslationUnit() (*ast.TranslationUnit, error) {
	topLevelNodes, err := parser.parseTopLevelNodes()
	if err != nil {
		return nil, err
	}
	return ast.NewTranslationUnit(parser.unitName, parser.rootScope, topLevelNodes), nil
}

func Parse(unitName string, scanner *scanner.Scanner, recorder *diagnostic.Recorder) (*ast.TranslationUnit, error) {
	parser := NewParser(unitName, scanner, recorder)
	return parser.ParseTranslationUnit()
}
