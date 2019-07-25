package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/code"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

type Factory struct {
	TokenReader token.Reader
	UnitName    string
	Recorder    *diagnostic.Recorder
}

func NewDefaultFactory() *Factory {
	return &Factory{
		UnitName: "undefined",
		Recorder: diagnostic.NewRecorder(),
	}
}

func (factory *Factory) WithTokenReader(reader token.Reader) *Factory {
	factory.TokenReader = reader
	return factory
}

func (factory *Factory) WithRecorder(recorder *diagnostic.Recorder) *Factory {
	factory.Recorder = recorder
	return factory
}

// NewParser creates a parser instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parser {
	parser := &Parser{
		rootScope:   code.NewRootScope(),
		tokenReader: factory.TokenReader,
		recorder:    factory.Recorder,
		unitName:    factory.UnitName,
	}
	parser.openBlock(token.NoIndent)
	parser.advance()
	return parser
}
