package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/scope"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

type Factory struct {
	TokenReader token.Reader
	UnitName    string
	Bag         *diagnostic.Bag
}

func NewDefaultFactory() *Factory {
	return &Factory{
		UnitName: "undefined",
		Bag:      diagnostic.NewBag(),
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

func (factory *Factory) WithRecorder(recorder *diagnostic.Bag) *Factory {
	factory.Bag = recorder
	return factory
}

// NewParser creates a parsing instance that parses the tokens of the given
// token.Reader and uses the 'unit' as its ast-root node. Errors while parsing
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parsing {
	parser := &Parsing{
		rootScope:   scope.NewRootScope(),
		tokenReader: factory.TokenReader,
		recorder:    factory.Bag,
		unitName:    factory.UnitName,
	}
	parser.openBlock(token.NoIndent)
	parser.advance()
	return parser
}
