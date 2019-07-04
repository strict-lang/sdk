package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"strings"
	"testing"
)

func TestParser_ParseTypeName(test *testing.T) {
	entries := []string{
		"number",
		"string",
		"list<number>",
		"list<list<number>>",
		"list < number >",
	}

	for _, entry := range entries {
		parser := NewTestParser(scanner.NewStringScanner(entry))
		name, err := parser.ParseTypeName()
		if err != nil {
			test.Errorf("unexpected error while parsing %s: %s", entry, err)
			continue
		}
		if name.FullName() != strings.ReplaceAll(entry, " ","") {
			test.Errorf("unexpected name %s, expected %s", name.FullName(), entry)
		}
	}
}
