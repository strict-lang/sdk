package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/source"
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
func (scanner *Scanner) GatherIdentifier() (string, error) {
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
