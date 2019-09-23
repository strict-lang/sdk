package scanning

import (
	"errors"
	"gitlab.com/strict-lang/sdk/compilation/token"
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
		next := scanning.char()
		if next == '"' {
			break
		}
		if next == '\n' {
			return "", ErrStringContainsLineFeed
		}
		if next == '\\' {
			_, ok := findEscapedCharacter(scanning.char())
			if !ok {
				return "", ErrInvalidEscapedChar
			}
			// TODO: Change this after backend emits something else
			builder.WriteRune('\\')
			scanning.advance()
			continue
		}
		builder.WriteRune(rune(next))
		scanning.advance()
	}
	return builder.String(), nil
}

func (scanning *Scanning) ScanStringLiteral() token.Token {
	literal, err := scanning.gatherStringLiteral()
	position := scanning.currentPosition()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewStringLiteralToken(literal, position, scanning.indent)
}
