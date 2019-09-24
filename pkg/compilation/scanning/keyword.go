package scanning

import (
	"errors"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

var (
	// ErrNoSuchKeyword is returned by the ScanKeyword() method when
	// a scanned identifier could not be found in the keyword-name-table.
	ErrNoSuchKeyword = errors.New("unknown keyword")
)

func (scanning *Scanning) ScanKeyword() token2.Token {
	keyword, err := scanning.gatherKeyword()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token2.NewKeywordToken(keyword, scanning.currentPosition(), scanning.indent)
}

// ScanKeyword scans a KeywordToken from the stream of characters.
func (scanning *Scanning) gatherKeyword() (token2.Keyword, error) {
	identifier, err := scanning.gatherIdentifier()
	if err != nil {
		return token2.InvalidKeyword, err
	}
	if keyword, ok := token2.KeywordByName(identifier); ok {
		return keyword, nil
	}
	return token2.InvalidKeyword, ErrNoSuchKeyword
}

func (scanning *Scanning) createInvalidKeyword(text string) token2.Token {
	return token2.NewInvalidToken(text, scanning.currentPosition(), scanning.indent)
}
