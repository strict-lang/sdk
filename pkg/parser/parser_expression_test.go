package parser

import (
	"com/github/BenjaminNitschke/Strict/pkg/ast"
	"com/github/BenjaminNitschke/Strict/pkg/scanner"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"1 + 2",
	}

	for _, entry := range entries {
		parser := createParser(input)
		expression, err := parser.ParseBinaryExpression(entry)
		if err != nil {
			test.Failf("unexpected error: %s", err.Error())
			continue
		}
		// TODO(merlinoayimwen): Verify result node
	}
}

func createParser(input string) *Parser {
	scanner := scanner.NewStringScanner(input)
	scope := ast.NewRootScope()
	// TODO(merlinosayimwen): Use diagnostics.Recorder
	unit := ast.NewTranslationUnit("test", scope, nil)
	return NewParser(unit, scanner, nil)
}
