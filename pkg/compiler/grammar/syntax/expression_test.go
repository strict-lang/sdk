package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"1 * 1",
		"1 + 1 * 2",
		"printf(\"%d\", limit(10) + 1))",
		"call(call(arg))",
		"1 isnt 1",
		"1 is 1 or 1 isnt 2",
		"random % 2 is 1",
		"index % 3 is 0 or index % 5 is 0",
		"!1",
		"(1 + 2)",
		"1 + 2 + 3",
		"(1 + 2) + 3",
	}

	for _, entry := range entries {
		testParsingBinaryExpression(test, entry)
	}
}

func testParsingBinaryExpression(test *testing.T, entry string) {
	parser := NewTestParser(lexical.NewStringScanning(entry))
	output, err := parser.parseExpression()
	if err != nil {
		test.Errorf("unexpected error while grammar (%s): %s", entry, err.Error())
		return
	}
	pretty.PrintColored(output)
}

type testEntry struct {
	code     string
	expected tree.Node
}

func testBinaryExpressionParsing(test *testing.T) {

}
