package parsing

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
	"gitlab.com/strict-lang/sdk/pkg/compilation/token"
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
import Neuron

Neuron[] inputNeurons
Neuron outputNeuron

method float Calculate(float left, float right)
  // float combinedWeights = combineWeights(left, right)
  if combinedWeights >= outputNeuron.Bias do
    return 1
  else
    return 0

method float combineWeights(float left, float right)
  float leftWeight = left * inputs[0].Weight
  float rightWeight = right * inputs[1].Weight
  return leftWeight + rightWeight
`
	tokens := scanning.NewStringScanning(entry)
	parser, bag := NewTestParserAndDiagnosticBag(tokens)
	_, err := parser.ParseTranslationUnit()
	if err != nil {
		test.Error(err)
	}
	diagnostics := bag.CreateDiagnostics(tokens.NewLineMap().PositionAtOffset)
	diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
}
