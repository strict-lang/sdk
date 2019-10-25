package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"strings"
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

type testParsingFunction = func(*Parsing) tree.Node

func ExpectResult(
	testing *testing.T,
	input string,
	expected tree.Node,
	parsingFunction testParsingFunction) {

	tokens := lexical.NewStringScanning(input)
	parser, bag := NewTestParserAndDiagnosticBag(tokens)
	result := parsingFunction(parser)
	defer func() {
		if failure := recover(); failure != nil {
			diagnostics := bag.CreateDiagnostics(tokens.NewLineMap().PositionAtOffset)
			diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
			testing.Errorf("Could not parse test input: %s", failure)
		}
	}()
	if result == nil {
		testing.Error("Result is nil")
		return
	}
	if !result.Matches(expected) {
		testing.Error(createUnexpectedTreeErrorMessage(input, expected, result))
	}
}

const testFormatIndent = 2
const testFormatLineBegin = "| "

func createUnexpectedTreeErrorMessage(
	input string, expected tree.Node, gotten tree.Node) string {

	formattedExpected := formatAndIndent(expected, testFormatIndent)
	formattedGotten := formatAndIndent(gotten, testFormatIndent)
	formattedInput := indentString(input, testFormatIndent, testFormatLineBegin)
	return fmt.Sprintf(`
The parser failed to parse following test-input correctly:
%s
Following tree was expected:
%s
Following tree was parsed:
%s
`, formattedInput, formattedExpected, formattedGotten)
}

func formatAndIndent(node tree.Node, indent int) string {
	formatted := pretty.Format(node)
	return indentString(formatted, indent, testFormatLineBegin)
}

func indentString(text string, indent int, lineBegin string) string {
	buffer := strings.Builder{}
	writeIndent(&buffer, indent)
	buffer.WriteString(lineBegin)
	for _, character := range text {
		buffer.WriteRune(character)
		if character == '\n' {
			writeIndent(&buffer, indent)
			buffer.WriteString(lineBegin)
		}
	}
	return buffer.String()
}

func writeIndent(buffer *strings.Builder, indent int) {
	for count := 0; count < indent; count++ {
		buffer.WriteRune(' ')
	}
}
