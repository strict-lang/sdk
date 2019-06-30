package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/source/linemap"
	"github.com/BenjaminNitschke/Strict/pkg/token"
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


func (scanner *Scanner) Pull() token.Token {
	return token.EndOfFile
}

func (scanner *Scanner) Peek() token.Token {
	return token.EndOfFile
}

func (scanner *Scanner) Last() token.Token {
	return token.EndOfFile
}