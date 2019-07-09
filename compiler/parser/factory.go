package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/scope"
	"gitlab.com/strict-lang/sdk/compiler/source/linemap"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

type Factory struct {
	TokenReader token.Reader
	Linemap *linemap.Linemap
	UnitName string
	Recorder *diagnostic.Recorder
}

func NewDefaultFactory() *Factory {
	return &Factory{
		UnitName: "undefined",
		Recorder: diagnostic.NewRecorder(),
		Linemap: linemap.NewBuilder().NewLinemap(),
	}
}

func (factory *Factory) WithTokenReader(reader token.Reader) *Factory {
	factory.TokenReader = reader
	return factory
}

// NewParser creates a parser instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parser {
	parser := &Parser{
		rootScope:   scope.NewRoot(),
		tokenReader: factory.TokenReader,
		recorder:    factory.Recorder,
		unitName:    factory.UnitName,
		linemap: 		 factory.Linemap,
	}
	parser.openBlock(token.NoIndent)
	parser.advance()
	return parser
}