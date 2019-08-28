package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/scope"
	"gitlab.com/strict-lang/sdk/compilation/token"
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

func (factory *Factory) WithUnitName(name string) *Factory {
	factory.UnitName = name
	return factory
}

func (factory *Factory) WithTokenReader(reader token.Reader) *Factory {
	factory.TokenReader = reader
	return factory
}

func (factory *Factory) WithRecorder(recorder *diagnostic.Recorder) *Factory {
	factory.Recorder = recorder
	return factory
}

// NewParser creates a parsing instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parsing {
	parser := &Parsing{
		rootScope:   scope.NewRootScope(),
		tokenReader: factory.TokenReader,
		recorder:    factory.Recorder,
		unitName:    factory.UnitName,
	}
	parser.openBlock(token.NoIndent)
	parser.advance()
	return parser
}
