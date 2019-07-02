package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

var (
	// ErrNoSuchKeyword is returned by the ScanKeyword() method when
	// a scanned identifier could not be found in the keyword-name-table.
	ErrNoSuchKeyword = errors.New("unknown keyword")
)

func (scanner *Scanner) ScanKeyword() token.Token {
	keyword, err := scanner.gatherKeyword()
	if err != nil {
		scanner.reportError(err)
		return scanner.createInvalidToken()
	}
	return token.NewKeywordToken(keyword, scanner.currentPosition(), scanner.indent)
}

// ScanKeyword scans a KeywordToken from the stream of characters.
func (scanner *Scanner) gatherKeyword() (token.Keyword, error) {
	identifier, err := scanner.gatherIdentifier()
	if err != nil {
		return token.InvalidKeyword, err
	}
	if keyword, ok := token.KeywordByName(identifier); ok {
		return keyword, nil
	}
	return token.InvalidKeyword, ErrNoSuchKeyword
}

func (scanner *Scanner) createInvalidKeyword(text string) token.Token {
	return token.NewInvalidToken(text, scanner.currentPosition(), scanner.indent)
}
