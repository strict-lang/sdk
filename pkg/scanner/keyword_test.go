package scanner

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"strings"
	"testing"
)

// TestGatherKnownKeywords passes all known keywords to a Scanner and ensures
// that its GatherKeyword method scans them correctly. If it fails to scan one
// of the keywords, the test fails.
func TestGatherKnownKeywords(test *testing.T) {
	for entryName, entry := range token.KeywordNames() {
		scanner := NewStringScanner(entryName)
		scanned, err := scanner.GatherKeyword()
		if err != nil {
			if err == ErrNoSuchKeyword {
				test.Errorf("scanner did not recognize keyword %s", entryName)
				return
			}
			test.Errorf("unexpected error %s", err.Error())
			continue
		}
		if scanned != entry {
			nameOfScanned := token.NameOfKind(entry)
			test.Errorf("scanned keyword %s, expected %s", nameOfScanned, entryName)
			return
		}
	}
}

func TestGatherInvalidKeywords(test *testing.T) {
	var entries = []string{
		"not_a_keyword",
		"AlsoNoKeyword",
	}

	for _, name := range token.KeywordNames() {
		// Appends the name of all known keywords as an upper case string to
		// the set of invalid entries. The language is case sensitive and thus
		// should not scan those strings as keywords.
		entries = append(entries, strings.ToUpper(name))
		// Appends the name of all known keywords with a leading underscore
		// to the set of invalid entries. The underscore should not be ignored.
		entries = append(entries, fmt.Sprintf("_%s", name))
	}

	for _, entry := range entries {
		scanner := NewStringScanner(entry)
		_, err := scanner.GatherKeyword()
		if err == nil {
			test.Errorf("scanner scanned invalid keyword %s", entry)
			return
		}
		if err != ErrNoSuchKeyword {
			test.Errorf("unexpected error %s", err.Error())
		}
	}
}
