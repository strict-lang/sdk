package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

var (
	ErrNoSuchKeyword = errors.New("unknown keyword")
)

func (scanner *Scanner) GatherKeyword() (token.Kind, error) {
	identifier, err := scanner.GatherIdentifier()
	if err != nil {
		return token.Invalid, err
	}
	if len(identifier) < token.ShortestKeywordNameLength() {
		return token.Invalid, ErrNoSuchKeyword
	}
	keyword, ok := token.KeywordByName(identifier)
	if !ok {
		return token.Invalid, ErrNoSuchKeyword
	}
	return keyword, nil
}
