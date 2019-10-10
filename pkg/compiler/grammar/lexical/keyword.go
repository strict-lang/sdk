package lexical

import (
	"errors"
	 "gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
)

var (
	// ErrNoSuchKeyword is returned by the ScanKeyword() method when
	// a scanned identifier could not be found in the keyword-name-table.
	ErrNoSuchKeyword = errors.New("unknown keyword")
)

func (scanning *Scanning) ScanKeyword() token.Token {
	keyword, err := scanning.gatherKeyword()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewKeywordToken(keyword, scanning.currentPosition(), scanning.indent)
}

// ScanKeyword scans a KeywordToken from the stream of characters.
func (scanning *Scanning) gatherKeyword() (token.Keyword, error) {
	identifier, err := scanning.gatherIdentifier()
	if err != nil {
		return token.InvalidKeyword, err
	}
	if keyword, ok := token.KeywordByName(identifier); ok {
		return keyword, nil
	}
	return token.InvalidKeyword, ErrNoSuchKeyword
}

func (scanning *Scanning) createInvalidKeyword(text string) token.Token {
	return token.NewInvalidToken(text, scanning.currentPosition(), scanning.indent)
}
