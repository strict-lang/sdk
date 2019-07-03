package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		`call(call(1))`,
		`1 + 1`,
		`!1`,
		`1 + 2 + 3`,
		`(1 + 2) + 3`,
		`printf("%d", limit(10) + 1)`,
	}

	for _, entry := range entries {
		testParsingBinaryExpression(test, entry)
	}
}

func testParsingBinaryExpression(test *testing.T, entry string) {
	parser := createParser(entry)
	defer parser.recorder.PrintAllEntries(diagnostic.NewTestPrinter(test))
	expression, err := parser.ParseExpression()
	if err != nil {
		test.Errorf("unexpected error while parsing (%s): %s", entry, err.Error())
		return
	}
	test.Log(expression)
}

func createParser(input string) *Parser {
	return NewParser("test", scanner.NewStringScanner(input), diagnostic.NewRecorder())
}
