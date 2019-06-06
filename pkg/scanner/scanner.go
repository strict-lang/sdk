package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/source/linemap"
	"strings"
)

type Scanner struct {
	reader         source.Reader
	linemapBuilder *linemap.Builder
}

func NewScanner(reader source.Reader) *Scanner {
	return &Scanner{
		reader:         reader,
		linemapBuilder: linemap.NewBuilder(),
	}
}

func NewStringScanner(input string) *Scanner {
	return NewScanner(source.NewStringReader(input))
}

type charMatcher func(source.Char) bool

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
