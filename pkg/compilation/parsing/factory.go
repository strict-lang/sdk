package parsing

import (
	code2 "gitlab.com/strict-lang/sdk/pkg/compilation/code"
	diagnostic2 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

// Factory is responsible for creating new parsing instances.
type Factory struct {
	tokens   token2.Stream
	unitName string
	bag      *diagnostic2.Bag
}

// NewDefaultFactory creates a factory with default values.
func NewDefaultFactory() *Factory {
	return &Factory{
		unitName: "undefined",
		bag:      diagnostic2.NewBag(),
	}
}

// WithUnitName sets the name of the translation unit.
func (factory *Factory) WithUnitName(name string) *Factory {
	factory.unitName = name
	return factory
}

// WithTokenStream set the source of tokens. This field is not copied per
// parser thus, creating multiple Parsings from a factory is not possible,
// unless the stream is changed each time.
func (factory *Factory) WithTokenStream(reader token2.Stream) *Factory {
	factory.tokens = reader
	return factory
}

// WithDiagnosticBag sets the diagnostic.Bag that diagnostics are reported to.
func (factory *Factory) WithDiagnosticBag(recorder *diagnostic2.Bag) *Factory {
	factory.bag = recorder
	return factory
}

// NewParser creates a parsing instance that parses the tokens of the given
// token.Stream and uses the 'unit' as its syntaxtree-root node. Errors while parsing
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parsing {
	parser := &Parsing{
		rootScope:   code2.NewRootScope(),
		tokenReader: factory.tokens,
		recorder:    factory.bag,
		unitName:    factory.unitName,
	}
	parser.openBlock(token2.NoIndent)
	parser.advance()
	parser.isAtBeginOfStatement = true
	return parser
}
