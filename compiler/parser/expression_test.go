package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/scope"
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
	// TODO(merlinosayimwen): Use diagnostics.Recorder
	unit := ast.NewTranslationUnit("test", scope.NewRoot(), []ast.Node{})
	return NewParser(&unit, scanner.NewStringScanner(input), diagnostic.NewRecorder())
}
