package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

var (
	// ErrInvalidIdentifier is returned when an identifier can not be scanned.
	ErrInvalidIdentifier = errors.New("invalid identifier")
)

func isIdentifierBegin(char source.Char) bool {
	return char.IsAlphabetic() || char == '_'
}

func isIdentifierChar(char source.Char) bool {
	return isIdentifierBegin(char) || char.IsNumeric()
}

// GatherIdentifier scans an identifier and returns it or an error if it fails.
func (scanner *Scanner) gatherIdentifier() (string, error) {
	leading, ok := scanner.scanMatching(isIdentifierBegin)
	if !ok {
		return "", ErrInvalidIdentifier
	}
	remaining, ok := scanner.scanAllMatching(isIdentifierChar)
	if !ok {
		return string(leading), nil
	}
	return string(leading) + remaining, nil
}

// tryToScanIdentifier tries to scan an identifier and records an error, if it fails.
func (scanner *Scanner) ScanIdentifier() token.Token {
	identifier, err := scanner.gatherIdentifier()
	position := scanner.currentPosition()
	if err != nil {
		scanner.reportError(err)
		return token.NewInvalidToken(scanner.reader.String(), position)
	}
	return token.NewIdentifierToken(identifier, position)
}
