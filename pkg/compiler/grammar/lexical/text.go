package lexical

import (
	"errors"
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/input"
	"strings"
)

const textCharacterLimit = 1024

var (
	errNoLeadingQuoteInString = errors.New("string literals does not begin with a quote")
	errStringContainsLineFeed = errors.New("string literal contains linefeed")
	errInvalidEscapedChar     = errors.New("literal contains invalid escaped char")
)

var escapedCharacters = map[rune]rune{
	't':  '\t',
	'n':  '\n',
	'f':  '\f',
	'r':  '\r',
	'b':  '\b',
	'\'': '\'',
	'"':  '"',
	'\\': '\\',
	'0':  rune(0),
}

func findEscapedCharacter(char input.Char) (input.Char, bool) {
	if escaped, ok := escapedCharacters[rune(char)]; ok {
		return input.Char(escaped), true
	}
	return input.EndOfFile, false
}

func (scanning *Scanning) gatherStringLiteral() (string, error) {
	if !scanning.tryToSkip('"') {
		return "", errNoLeadingQuoteInString
	}
	var builder strings.Builder
	for count := 0; count < textCharacterLimit; count++ {
		switch next := scanning.char(); next {
		case '"':
			scanning.advance()
			return builder.String(), nil
		case '\n':
			return "", errStringContainsLineFeed
		case '\\':
			if _, ok := findEscapedCharacter(scanning.peekChar()); !ok {
				return "", errInvalidEscapedChar
			}
			builder.WriteRune('\\')
			scanning.advance()
		}
		builder.WriteRune(rune(scanning.char()))
		scanning.advance()
	}
	return builder.String(), nil
}

func (scanning *Scanning) scanStringLiteral() token.Token {
	literal, err := scanning.gatherStringLiteral()
	position := scanning.currentPosition()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewStringLiteralToken(literal, position, scanning.indent)
}
