package scanning

import (
	"errors"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"strings"
)

const textCharacterLimit = 1024

var (
	ErrNoLeadingQuoteInString = errors.New("string literals does not begin with a quote")
	ErrStringContainsLineFeed = errors.New("string literal contains linefeed")
	ErrInvalidEscapedChar     = errors.New("literal contains invalid escaped char")
)

func (scanning *Scanning) gatherStringLiteral() (string, error) {
	if !scanning.tryToSkip('"') {
		return "", ErrNoLeadingQuoteInString
	}
	var builder strings.Builder
	for count := 0; count < textCharacterLimit; count++ {
		switch next := scanning.char(); next {
		case '"':
			scanning.advance()
			return builder.String(), nil
		case '\n':
			return "", ErrStringContainsLineFeed
		case '\\':
			if _, ok := findEscapedCharacter(scanning.peekChar()); !ok {
				return "", ErrInvalidEscapedChar
			}
			builder.WriteRune('\\')
			scanning.advance()
		}
		builder.WriteRune(rune(scanning.char()))
		scanning.advance()
	}
	return builder.String(), nil
}

func (scanning *Scanning) ScanStringLiteral() token2.Token {
	literal, err := scanning.gatherStringLiteral()
	position := scanning.currentPosition()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token2.NewStringLiteralToken(literal, position, scanning.indent)
}
