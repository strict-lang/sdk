package scanning

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"strings"
)

type charMatcher func(source.Char) bool

func (scanning *Scanning) offset() source.Offset {
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

func (scanning *Scanning) scanMatching(matcher charMatcher) (source.Char, bool) {
	char := scanning.char()
	if matcher(scanning.char()) {
		scanning.advance()
		return char, true
	}
	return source.EndOfFile, false
}

// tryToSkip consumes the next character if it has the same id, as the one
// passed to the function, otherwise the index remains the same.
func (scanning *Scanning) tryToSkip(char source.Char) bool {
	next := scanning.char()
	if next != char {
		return false
	}
	scanning.advance()
	return true
}

func (scanning *Scanning) tryToSkipMultiple(char source.Char, amount int) bool {
	for count := 0; count < amount; count++ {
		if !scanning.tryToSkip(char) {
			return false
		}
	}
	return true
}

func (scanning *Scanning) createPositionToOffset(begin source.Offset) token.Position {
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
