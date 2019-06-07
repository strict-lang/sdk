package scanner

import "testing"

func TestScanningStringLiteralBody(test *testing.T) {
	entries := []string{
		"\"This is the content of a string\"",
		"\"3123{}[]<>\"",
		"\"\"", // empty string
	}
	for _, entry := range entries {
		scanner := NewStringScanner(entry)
		token, err := scanner.ScanStringLiteral()
		if err != nil {
			test.Errorf("unexpected error while scanning: %s", entry)
			continue
		}
		entryBody := removeSurroundingQuotes(entry)
		if token.Value != entryBody {
			test.Errorf("scanned '%s' but expected '%s'", token.Value, entryBody)
		}
	}
}

// Removes the first and last character from the string.
func removeSurroundingQuotes(literal string) string {
	return literal[1 : len(literal)-1]
}

func TestScanningInvalidStringLiteralBody(test *testing.T) {
	entries := map[string]error{
		"\"This contains a \nLinebreak\"": ErrStringContainsLineFeed,
		"\"\\p Invalid escaped char \"":   ErrInvalidEscapedChar,
	}
	for entry, expected := range entries {
		scanner := NewStringScanner(entry)
		token, err := scanner.ScanStringLiteral()
		if err == nil {
			test.Errorf("scanned invalid string %s as %s", entry, token.Value)
			continue
		}
		if err != expected {
			test.Errorf("unexpected error %s, expected %s", entry, expected)
		}
	}
}
