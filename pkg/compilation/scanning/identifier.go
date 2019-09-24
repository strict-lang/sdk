package scanning

import (
	"errors"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

var (
	// ErrInvalidIdentifier is returned when an identifier can not be scanned.
	ErrInvalidIdentifier = errors.New("invalid identifier")
)

func isIdentifierBegin(char source2.Char) bool {
	return char.IsAlphabetic() || char == '_'
}

func isIdentifierChar(char source2.Char) bool {
	return isIdentifierBegin(char) || char.IsNumeric()
}

// GatherIdentifier scans an identifier and returns it or an error if it fails.
func (scanning *Scanning) gatherIdentifier() (string, error) {
	leading, ok := scanning.scanMatching(isIdentifierBegin)
	if !ok {
		return "", ErrInvalidIdentifier
	}
	remaining, ok := scanning.scanAllMatching(isIdentifierChar)
	if !ok {
		return string(leading), nil
	}
	return string(leading) + remaining, nil
}

// tryToScanIdentifier tries to scanning an identifier and records an error, if it fails.
func (scanning *Scanning) ScanIdentifier() token2.Token {
	identifier, err := scanning.gatherIdentifier()
	position := scanning.currentPosition()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token2.NewIdentifierToken(identifier, position, scanning.indent)
}

func (scanning *Scanning) ScanIdentifierOrKeyword() token2.Token {
	identifier, err := scanning.gatherIdentifier()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	if keyword, ok := token2.KeywordByName(identifier); ok {
		return token2.NewKeywordToken(keyword, scanning.currentPosition(), scanning.indent)
	}
	return token2.NewIdentifierToken(identifier, scanning.currentPosition(), scanning.indent)
}
