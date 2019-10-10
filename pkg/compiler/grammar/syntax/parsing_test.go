package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"testing"
)

func NewTestParser(tokens token.Stream) *Parsing {
	return NewDefaultFactory().WithTokenStream(tokens).NewParser()
}

func NewTestParserAndDiagnosticBag(tokens token.Stream) (*Parsing, *diagnostic.Bag) {
	bag := diagnostic.NewBag()
	return NewDefaultFactory().
		WithDiagnosticBag(bag).
		WithTokenStream(tokens).
		NewParser(), bag
}

func TestParseTopLevelStatements(test *testing.T) {
	const entry = `
method decreaseWeights()
  inputNeurons[0].Weights -= trainingLambda
  inputNeurons[1].Weights -= trainingLambda
`
	tokens := lexical.NewStringScanning(entry)
	parser, bag := NewTestParserAndDiagnosticBag(tokens)
	unit, err := parser.ParseTranslationUnit()
	if err != nil {
		test.Error(err)
	}
	diagnostics := bag.CreateDiagnostics(tokens.NewLineMap().PositionAtOffset)
	diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
	pretty.PrintColored(unit)
}
