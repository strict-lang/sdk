package syntax

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
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

type testParsingFunction = func(*Parsing) tree.Node

type ParserTestEntry struct {
	Input          string
	ExpectedOutput tree.Node
}

// ExpectAllResults tests multiple entries using the same function.
func ExpectAllResults(
	testing *testing.T,
	entries []ParserTestEntry,
	parsingFunction testParsingFunction) {

	for _, entry := range entries {
		ExpectResult(testing, entry.Input, entry.ExpectedOutput, parsingFunction)
	}
}

func ExpectError(
	testing *testing.T,
	input string,
	parsingFunction testParsingFunction,
	matcher func(err error) bool) {

	// TODO: Fix this combined with the diagnostics
	tokens := lexical.NewStringScanning(input)
	parser := NewTestParser(tokens)
	defer func() {
		if failure := recover(); failure != nil {
			err := extractErrorFromPanic(failure)
			if !matcher(err) {
				// testing.Errorf("unexpected error: %s", err)
			}
		} else {
			// testing.parsingError("no error was reported")
		}
	}()
	parsingFunction(parser)
}

func ExpectResult(
	testing *testing.T,
	input string,
	expected tree.Node,
	parsingFunction testParsingFunction) {

	tokens := lexical.NewStringScanning(input)
	parser, bag := NewTestParserAndDiagnosticBag(tokens)
	defer func() {
		if failure := recover(); failure != nil {
			diagnostics := bag.CreateDiagnostics(tokens.NewLineMap().PositionAtOffset)
			diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
			testing.Errorf("%s", failure)
		}
	}()
	result := parsingFunction(parser)
	if result == nil {
		testing.Error("Result is nil")
		return
	}
	if !result.Matches(expected) {
		testing.Error(createUnexpectedTreeErrorMessage(input, expected, result))
	}
	expectEmptyStructureStack(testing, parser)
}

func expectEmptyStructureStack(testing *testing.T, parsing *Parsing) {
	if parsing.structureStack.isEmpty() {
		return
	}
	message := strings.Builder{}
	message.WriteString(`
There are remaining structures on the parsings structure-stack.
This indicates that the parser has begun to parse one or more structures, that have not
been completed. This should not happen and is the result of a bug/missing completion in the code.
Following structures were added but never removed from the stack (starting at the top element):

`)
	for _, remaining := range parsing.structureStack.listRemainingElementsOrdered() {
		item := fmt.Sprintf(" - %s\n", remaining.nodeKind)
		message.WriteString(item)
	}
	message.WriteString(`
Following list shows the structure-stack's history. It contains all
elements that have been pushed to and popped from the stack:

`)
	message.WriteString(parsing.structureStack.createHistoryDump())
	testing.Error(message.String())
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
