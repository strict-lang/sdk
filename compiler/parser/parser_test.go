package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

func NewTestParser(tokens token.Reader) *Parser {
	return NewParser("test", tokens, diagnostic.NewRecorder())
}
