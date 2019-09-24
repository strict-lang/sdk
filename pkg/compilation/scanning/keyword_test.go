package scanning

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"strings"
	"testing"
)

// TestGatherKnownKeywords passes all known keywords to a Scanning and ensures
// that its GatherKeyword method scans them correctly. If it fails to scanning one
// of the keywords, the test fails.
func TestGatherKnownKeywords(test *testing.T) {
	for entryName, entry := range token.KeywordNames() {
		scanner := NewStringScanning(entryName)
		scanned, err := scanner.gatherKeyword()
		if err != nil {
			if err == ErrNoSuchKeyword {
				test.Errorf("scanning did not recognize keyword %s", entryName)
				return
			}
			test.Errorf("unexpected error %s", err.Error())
			continue
		}
		if scanned != entry {
			test.Errorf("scanned keyword %s, expected %s", scanned.String(), entryName)
			return
		}
	}
}

func TestGatherInvalidKeywords(test *testing.T) {
	var entries = []string{
		"not_a_keyword",
		"AlsoNoKeyword",
	}

	for name := range token.KeywordNames() {
		// Appends the name of all known keywords as an upper case string to
		// the set of invalid entries. The language is case sensitive and thus
		// should not scanning those strings as keywords.
		entries = append(entries, strings.ToUpper(name))
		// Appends the name of all known keywords with a leading underscore
		// to the set of invalid entries. The underscore should not be ignored.
		entries = append(entries, fmt.Sprintf("_%s", name))
	}

	for _, entry := range entries {
		scanner := NewStringScanning(entry)
		_, err := scanner.gatherKeyword()
		if err == nil {
			test.Errorf("scanning scanned invalid keyword %s", entry)
			return
		}
		if err != ErrNoSuchKeyword {
			test.Errorf("unexpected error %s", err.Error())
		}
	}
}
