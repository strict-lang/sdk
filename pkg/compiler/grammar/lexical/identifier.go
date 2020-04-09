package lexical

import (
	"errors"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/input"
)

var (
	// errInvalidIdentifier is returned when an identifier can not be scanned.
	errInvalidIdentifier = errors.New("invalid identifier")
)

func isIdentifierBegin(char input.Char) bool {
	return char.IsAlphabetic() || char == '_'
}

func isIdentifierChar(char input.Char) bool {
	return isIdentifierBegin(char) || char.IsNumeric()
}

// GatherIdentifier scans an identifier and returns it or an error if it fails.
func (scanning *Scanning) gatherIdentifier() (string, error) {
	leading, ok := scanning.scanMatching(isIdentifierBegin)
	if !ok {
		return "", errInvalidIdentifier
	}
	remaining, ok := scanning.scanAllMatching(isIdentifierChar)
	if !ok {
		return string(leading), nil
	}
	return string(leading) + remaining, nil
}

// scanIdentifier tries to scanning an identifier and records an error, if it fails.
func (scanning *Scanning) scanIdentifier() token.Token {
	identifier, err := scanning.gatherIdentifier()
	position := scanning.currentPosition()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewIdentifierToken(identifier, position, scanning.indent)
}

func (scanning *Scanning) scanIdentifierOrKeyword() token.Token {
	identifier, err := scanning.gatherIdentifier()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	if keyword, ok := token.KeywordByName(identifier); ok {
		return token.NewKeywordToken(keyword, scanning.currentPosition(), scanning.indent)
	}
	return token.NewIdentifierToken(identifier, scanning.currentPosition(), scanning.indent)
}
