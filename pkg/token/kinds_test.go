package token

import (
	"testing"
)

func TestKeywordsAreKeywords(test *testing.T) {
	for keyword := keywordsBegin + 1; keyword < keywordsEnd; keyword++ {
		if !keyword.IsKeyword() {
			name := keyword.Name()
			test.Errorf("%s is not considered a keyword", name)
		}
	}
}
