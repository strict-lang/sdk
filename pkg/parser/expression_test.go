package parser

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/scanner"
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"a + b",
	}

	for _, entry := range entries {
		parser := createParser(entry)
		// TODO(merlinosayimwen): Verify result node
		expression, err := parser.ParseBinaryExpression()
		if err != nil {
			test.Errorf("unexpected error: %s", err.Error())
			continue
		}
		fmt.Println(expression)
	}
}

func createParser(input string) *Parser {
	// TODO(merlinosayimwen): Use diagnostics.Recorder
	unit := ast.NewTranslationUnit("test", scope.NewRoot(), []ast.Node{})
	return NewParser(&unit, scanner.NewStringScanner(input), diagnostic.NewRecorder())
}
