package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/compiler/source"
	"github.com/BenjaminNitschke/Strict/compiler/token"
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
		return scanner.createInvalidToken()
	}
	return token.NewIdentifierToken(identifier, position, scanner.indent)
}

func (scanner *Scanner) ScanIdentifierOrKeyword() token.Token {
	identifier, err := scanner.gatherIdentifier()
	if err != nil {
		scanner.reportError(err)
		return scanner.createInvalidToken()
	}
	if keyword, ok := token.KeywordByName(identifier); ok {
		return token.NewKeywordToken(keyword, scanner.currentPosition(), scanner.indent)
	}
	return token.NewIdentifierToken(identifier, scanner.currentPosition(), scanner.indent)
}