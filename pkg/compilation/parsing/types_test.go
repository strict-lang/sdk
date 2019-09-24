package parsing

import (
	scanning2 "gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
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
		parser := NewTestParser(scanning2.NewStringScanning(entry))
		name, err := parser.parseTypeName()
		if err != nil {
			test.Errorf("unexpected error while parsing %s: %s", entry, err)
			continue
		}
		if name.FullName() != strings.Replace(entry, " ", "", 10) {
			test.Errorf("unexpected name %s, expected %s", name.FullName(), entry)
		}
	}
}