package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

var (
	// ErrNoSuchKeyword is returned by the ScanKeyword() method when
	// a scanned identifier could not be found in the keyword-name-table.
	ErrNoSuchKeyword = errors.New("unknown keyword")
)

// ScanKeyword scans a KeywordToken from the stream of characters.
func (scanner *Scanner) ScanKeyword() (token.Token, error) {
	beginOffset := scanner.offset()
	identifier, err := scanner.gatherIdentifier()
	if err != nil {
		return scanner.createInvalidKeyword(beginOffset, identifier), err
	}

	keyword, ok := token.KeywordByName(identifier)
	if !ok {
		invalid := scanner.createInvalidKeyword(beginOffset, identifier)
		return invalid, ErrNoSuchKeyword
	}

	position := scanner.createPositionToOffset(beginOffset)
	return token.NewKeywordToken(keyword, position), nil
}

func (scanner *Scanner) createInvalidKeyword(
	beginOffset source.Offset, text string) token.Token {

	position := scanner.createPositionToOffset(beginOffset)
	return token.NewInvalidToken(text, position)
}
