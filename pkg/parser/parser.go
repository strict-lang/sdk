package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

// Parser parses an AST from a stream of tokens.
type Parser struct {
	tokens   token.Reader
	unit     *ast.TranslationUnit
	recorder *diagnostics.Recorder
}

// NewParser creates a parser instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func NewParser(unit *ast.TranslationUnit, tokens token.Reader, recorder *diagnostics.Recorder) *Parser {
	return &Parser{
		unit:     unit,
		tokens:   tokens,
		recorder: recorder,
	}
}

func (parser *Parser) skipOperator(operator token.Operator) (bool, error) {
	if ok, err := parser.expectOperator(); !ok {
		return false, err
	}
	parser.tokens.Pull()
	return true, nil
}

// expectOperator peeks the next token and expects it to be the passed operator,
// else false is returned and an error is recorded.
func (parser *Parser) expectOperator(operator token.Operator) (bool, error) {
	peek := parser.tokens.Peek()
	if !peek.IsOperator() {
		return false
	}
	ok := peek.(Operator).Operator() == operator
	return ok, nil
}

func (parser *Parser) expectKeyword(keyword token.Keyword) bool {
	return false
}

func (parser *Parser) expectIdentifier() bool {
	return false
}

func (parser *Parser) isLookingAtKeyword(keyword token.Keyword) bool {
	return false
}
