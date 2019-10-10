package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/lexical"
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
		parser := NewTestParser(lexical.NewStringScanning(entry))
		name, err := parser.parseTypeName()
		if err != nil {
			test.Errorf("unexpected error while grammar %s: %s", entry, err)
			continue
		}
		if name.FullName() != strings.Replace(entry, " ", "", 10) {
			test.Errorf("unexpected name %s, expected %s", name.FullName(), entry)
		}
	}
}
