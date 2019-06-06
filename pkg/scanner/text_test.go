package scanner

import "testing"

func TestScanningStringLiteralBody(test *testing.T) {
	// entries := []string {
	// "This is the content of a string",
	// "3123{}[]<>",
	// }

}

func TestScanningTerminatingQuote(test *testing.T) {

}

func TestScanningEscapedCharactersInStringLiteral(test *testing.T) {

}

func TestScanningInvalidStringLiteralBody(test *testing.T) {
	entries := map[string]error{
		"This contains a \nLinebreak": ErrStringContainsLineFeed,
		"\\p Invalid escaped char ":   ErrInvalidEscapedChar,
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
