package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"strings"
)

const textCharacterLimit = 1024

var (
	ErrNoLeadingQuoteInString = errors.New("string literals does not begin with a quote")
	ErrStringContainsLineFeed = errors.New("string literal contains linefeed")
	ErrInvalidEscapedChar     = errors.New("literal contains invalid escaped char")
)

func (scanner *Scanner) ScanStringLiteral() (token.Token, error) {
	if !scanner.tryToSkip('"') {
		return token.NewAnonymousInvalidToken(), ErrNoLeadingQuoteInString
	}
	var builder strings.Builder
	beginIndex := scanner.reader.Index()
	for count := 0; count < textCharacterLimit; count++ {
		next := scanner.reader.Pull()
		if next == '"' {
			break
		}
		if next == '\n' {
			return token.NewAnonymousInvalidToken(), ErrStringContainsLineFeed
		}
		if next == '\\' {
			escaped, ok := findEscapedCharacter(scanner.reader.Pull())
			if !ok {
				return token.NewAnonymousInvalidToken(), ErrInvalidEscapedChar
			}
			builder.WriteRune(rune(escaped))
			continue
		}
		builder.WriteRune(rune(next))
	}
	position := token.Position{
		Begin: beginIndex,
		End:   scanner.reader.Index(),
	}
	return token.NewStringLiteralToken(builder.String(), position), nil
}
