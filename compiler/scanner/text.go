package scanner

import (
	"errors"
	"gitlab.com/strict-lang/sdk/compiler/token"
	"strings"
)

const textCharacterLimit = 1024

var (
	ErrNoLeadingQuoteInString = errors.New("string literals does not begin with a quote")
	ErrStringContainsLineFeed = errors.New("string literal contains linefeed")
	ErrInvalidEscapedChar     = errors.New("literal contains invalid escaped char")
)

func (scanner *Scanner) gatherStringLiteral() (string, error) {
	if !scanner.tryToSkip('"') {
		return "", ErrNoLeadingQuoteInString
	}
	var builder strings.Builder
	for count := 0; count < textCharacterLimit; count++ {
		next := scanner.reader.Pull()
		if next == '"' {
			break
		}
		if next == '\n' {
			return "", ErrStringContainsLineFeed
		}
		if next == '\\' {
			escaped, ok := findEscapedCharacter(scanner.reader.Pull())
			if !ok {
				return "", ErrInvalidEscapedChar
			}
			builder.WriteRune(rune(escaped))
			continue
		}
		builder.WriteRune(rune(next))
	}
	return builder.String(), nil
}

func (scanner *Scanner) ScanStringLiteral() token.Token {
	literal, err := scanner.gatherStringLiteral()
	position := scanner.currentPosition()
	if err != nil {
		scanner.reportError(err)
		return scanner.createInvalidToken()
	}
	return token.NewStringLiteralToken(literal, position, scanner.indent)
}
