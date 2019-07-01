package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"strings"
)

type charMatcher func(source.Char) bool

func (scanner *Scanner) offset() source.Offset {
	return scanner.reader.Index()
}

func (scanner *Scanner) scanAllMatching(matcher charMatcher) (string, bool) {
	var builder strings.Builder
	for {
		if !matcher(scanner.reader.Peek()) {
			break
		}
		builder.WriteRune(rune(scanner.reader.Pull()))
	}
	return builder.String(), builder.Len() > 0
}

func (scanner *Scanner) scanMatching(matcher charMatcher) (source.Char, bool) {
	if !matcher(scanner.reader.Peek()) {
		return 0, false
	}
	return scanner.reader.Pull(), true
}

// tryToSkip consumes the next character if it has the same id, as the one
// passed to the function, otherwise the index remains the same.
func (scanner *Scanner) tryToSkip(char source.Char) bool {
	next := scanner.reader.Peek()
	if next != char {
		return false
	}
	scanner.reader.Pull()
	return true
}

func (scanner *Scanner) tryToSkipMultiple(char source.Char, amount int) bool {
	for count := 0; count < amount; count++ {
		if !scanner.tryToSkip(char) {
			return false
		}
	}
	return true
}

func (scanner *Scanner) createPositionToOffset(begin source.Offset) token.Position {
	return token.Position{
		Begin: begin,
		End:   scanner.offset(),
	}
}

func (scanner *Scanner) currentPosition() token.Position {
	return token.Position{
		Begin: scanner.reader.internalIndex,
		End:   scanner.offset(),
	}
}

func (scanner *Scanner) SkipWhitespaces() {
	for {
		next := scanner.reader.Pull();
		if next == '\n' {
			scanner.incrementLineIndex()
			continue
		}
		if next.IsWhitespace() || next == '\r' {
			continue
		}
		return
	}
}
