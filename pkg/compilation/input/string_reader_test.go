package input

import (
	"strings"
	"testing"
)

func TestRecreateInput(test *testing.T) {
	entries := []string{
		"Hello, World!",
		"someIdentifier",
		"123321234215123",
		"               ",
		")()()()()()()()(",
		"", // Empty
	}

	var builder strings.Builder
	for _, entry := range entries {
		reader := NewStringReader(entry)
		for {
			next := reader.Pull()
			if next == EndOfFile {
				break
			}
			builder.WriteRune(rune(next))
		}
		read := builder.String()
		if read != entry {
			test.Errorf("reader read '%s', expected '%s'", read, entry)
		}
		builder.Reset()
	}

}
