package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/scanner"
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"1 + 2",
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
	unit := ast.NewTranslationUnit("test", scope.NewRootScope(), []ast.Node{})
	return NewParser(&unit, scanner.NewStringScanner(input), diagnostic.NewRecorder())
}
