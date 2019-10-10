package lexical

import "testing"

// TestInvalidIdentifiers passes the Scanning some invalid identifiers to
// ensures that the GatherIdentifier method fails scanning them and returns
// an ErrInvalidIdentifier error. If one of the entries is scanned successfully
// or scanning it produces an unexpected error, the test fails.
func TestInvalidIdentifiers(test *testing.T) {
	var entries = []string{
		"0leadingNumber",
		"+leadingSymbol",
		"{curlyBracket}",
		"[blockBracket]",
		"(parentheses)",
		"", // empty strings
		"0",
	}

	for _, entry := range entries {
		scanner := NewStringScanning(entry)
		identifier, err := scanner.gatherIdentifier()
		if err == nil {
			test.Errorf("scanning scanned invalid identifier %s as %s", entry, identifier)
			continue
		}
		if err != ErrInvalidIdentifier {
			test.Fatalf("unexpected error %s", err.Error())
		}
	}
}

// TestValidIdentifiers passes the Scanning some valid identifiers to ensure that
// the GatherIdentifier method scans them without errors and produces a string
// that matches the entry. If it returns an error or fails to scanning the entry
// successfully, the test fails.
func TestValidIdentifiers(test *testing.T) {
	var entries = []string{
		"a",  // single letter
		"_",  // single underscore
		"__", // multiple underscores
		"PascalCase",
		"mixedCaseIdentifier",
		"lowercase",
		"lower_underscore",
		"UPPER_UNDERSCORE",
		"_leadingUnderscore",
		"trailingUnderscore_",
		"trailingNumber0",
		"m1x3d",
		"__more_underscores__",
	}

	for _, entry := range entries {
		scanner := NewStringScanning(entry)
		identifier, err := scanner.gatherIdentifier()
		if err != nil {
			test.Errorf("unexpected error %s", err.Error())
		}
		if identifier != entry {
			test.Errorf("scanning wrongly scanned entry %s as '%s'", entry, identifier)
		}
	}
}
