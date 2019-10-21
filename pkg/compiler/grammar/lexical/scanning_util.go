package lexical

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

type charMatcher func(input.Char) bool

func (scanning *Scanning) offset() input.Offset {
	return scanning.input.Index()
}

func (scanning *Scanning) scanAllMatching(matcher charMatcher) (string, bool) {
	var builder strings.Builder
	for {
		if scanning.input.IsExhausted() || !matcher(scanning.char()) {
			break
		}
		builder.WriteRune(rune(scanning.char()))
		scanning.advance()
	}
	return builder.String(), builder.Len() > 0
}

func (scanning *Scanning) scanMatching(matcher charMatcher) (input.Char, bool) {
	char := scanning.char()
	if matcher(scanning.char()) {
		scanning.advance()
		return char, true
	}
	return input.EndOfFile, false
}

// tryToSkip consumes the next character if it has the same id, as the one
// passed to the function, otherwise the index remains the same.
func (scanning *Scanning) tryToSkip(char input.Char) bool {
	next := scanning.char()
	if next != char {
		return false
	}
	scanning.advance()
	return true
}

func (scanning *Scanning) tryToSkipMultiple(char input.Char, amount int) bool {
	for count := 0; count < amount; count++ {
		if !scanning.tryToSkip(char) {
			return false
		}
	}
	return true
}

func (scanning *Scanning) createPositionToOffset(begin input.Offset) token.Position {
	return token.Position{
		BeginOffset: begin,
		EndOffset:   scanning.offset(),
	}
}

func (scanning *Scanning) currentPosition() token.Position {
	return token.Position{
		BeginOffset: scanning.input.Index(),
		EndOffset:   scanning.offset(),
	}
}

func (scanning *Scanning) skipWhitespaces() (token.Token, bool) {
	for {
		switch char := scanning.char(); char {
		case '\n':
			scanning.advance()
			if scanning.input.IsExhausted() {
				return scanning.endOfFile(), true
			}
			if endOfStatement, ok := scanning.incrementLineIndex(); ok {
				return endOfStatement, true
			}
			continue
		case ' ':
			scanning.addWhitespaceIndent()
			break
		case '\t':
			scanning.addTabIndent()
			break
		case '\r':
			break
		default:
			return nil, false
		}
		scanning.advance()
		if scanning.input.IsExhausted() {
			return nil, false
		}
	}
}

func (scanning *Scanning) addTabIndent() {
	scanning.addIndent(TabIndent)
}

func (scanning *Scanning) addWhitespaceIndent() {
	scanning.addIndent(WhitespaceIndent)
}

func (scanning *Scanning) addIndent(indent token.Indent) {
	if scanning.updateIndent {
		scanning.indent += indent
	}
}

func ScanAllTokens(scanner *Scanning) []token.Token {
	var tokens []token.Token
	for {
		next := scanner.Pull()
		if token.IsEndOfFileToken(next) {
			break
		}
		tokens = append(tokens, next)
	}
	return tokens
}
