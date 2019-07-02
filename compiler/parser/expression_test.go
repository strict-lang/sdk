package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"a + b",
	}

	for _, entry := range entries {
		parser := createParser(entry)
		// TODO(merlinosayimwen): Verify result node
		_, err := parser.ParseBinaryExpression()
		if err != nil {
			test.Errorf("unexpected error: %s", err.Error())
			continue
		}
	}
}

func createParser(input string) *Parser {
	return NewParser("test", scanner.NewStringScanner(input), diagnostic.NewRecorder())
}
