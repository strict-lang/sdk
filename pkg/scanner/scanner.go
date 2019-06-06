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
