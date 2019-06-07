package token

import "testing"

func TestKeywordByName(test *testing.T) {
	for keyword := keywordsBegin + 1; keyword < keywordsEnd; keyword++ {
		name := NameOfKind(keyword)
		found, ok := KeywordByName(name)
		if !ok {
			test.Errorf("keyword %s not found", name)
			return
		}
		if found != keyword {
			nameOfFound := NameOfKind(found)
			test.Errorf("expected kind %s but got %s", name, nameOfFound)
		}
	}
}

func TestKeywordByInvalidName(test *testing.T) {
	var entries = []string{
		"",
		"123",
		"not_a_keyword",
		"NeitherAKeyword",
	}

	for _, entry := range entries {
		found, ok := KeywordByName(entry)
		if !ok {
			continue
		}
		nameOfFound := NameOfKind(found)
		test.Errorf("expected nothing but got %s", nameOfFound)
	}
}
