package scanning

import (
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"strings"
)

type charMatcher func(source.Char) bool

func (scanning *Scanning) offset() source.Offset {
	return scanning.reader.Index()
}

func (scanning *Scanning) scanAllMatching(matcher charMatcher) (string, bool) {
	var builder strings.Builder
	for {
		if !matcher(scanning.reader.Peek()) {
			break
		}
		builder.WriteRune(rune(scanning.reader.Pull()))
	}
	return builder.String(), builder.Len() > 0
}

func (scanning *Scanning) scanMatching(matcher charMatcher) (source.Char, bool) {
	if !matcher(scanning.reader.Peek()) {
		return 0, false
	}
	return scanning.reader.Pull(), true
}

// tryToSkip consumes the next character if it has the same id, as the one
// passed to the function, otherwise the index remains the same.
func (scanning *Scanning) tryToSkip(char source.Char) bool {
	next := scanning.reader.Peek()
	if next != char {
		return false
	}
	scanning.reader.Pull()
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
		BeginOffset: scanning.reader.Index(),
		EndOffset:   scanning.offset(),
	}
}

func (scanning *Scanning) SkipWhitespaces() (token.Token, bool) {
	for {
		peek := scanning.reader.Peek()
		if peek == '\n' {
			scanning.reader.Pull()
			if endOfStatement, ok := scanning.incrementLineIndex(); ok {
				return endOfStatement, true
			}
			continue
		}
		if peek == ' ' {
			scanning.addWhitespaceIndent()
			scanning.reader.Pull()
			continue
		}
		if peek == '\t' {
			scanning.addTabIndent()
			scanning.reader.Pull()
			continue
		}
		if peek == '\r' {
			scanning.reader.Pull()
			continue
		}
		return nil, false
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
