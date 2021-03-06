package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
)

// Factory is responsible for creating new grammar instances.
type Factory struct {
	tokens   token.StreamWithLineMap
	unitName string
	bag      *diagnostic.Bag
}

// NewDefaultFactory creates a factory with default values.
func NewDefaultFactory() *Factory {
	return &Factory{
		unitName: "undefined",
		bag:      diagnostic.NewBag(),
	}
}

// WithUnitName sets the name of the translation unit.
func (factory *Factory) WithUnitName(name string) *Factory {
	factory.unitName = name
	return factory
}

// WithTokenStream set the input of tokens. This field is not copied per
// parser thus, creating multiple Parsings from a factory is not possible,
// unless the stream is changed each time.
func (factory *Factory) WithTokenStream(reader token.StreamWithLineMap) *Factory {
	factory.tokens = reader
	return factory
}

// WithDiagnosticBag sets the diagnostic.Bag that diagnostics are reported to.
func (factory *Factory) WithDiagnosticBag(recorder *diagnostic.Bag) *Factory {
	factory.bag = recorder
	return factory
}

// NewParser creates a grammar instance that parses the tokens of the given
// token.Stream and uses the 'unit' as its tree-root node. Errors while grammar
// are recorded by the 'recorder'.
func (factory *Factory) NewParser() *Parsing {
	parser := &Parsing{
		tokenReader:    factory.tokens,
		recorder:       factory.bag,
		unitName:       factory.unitName,
		structureStack: newRecordingStructureStack(),
		currentMethod:  notParsingMethod,
	}
	parser.openBlock(token.NoIndent)
	parser.advance()
	parser.statementBegin = true
	return parser
}
