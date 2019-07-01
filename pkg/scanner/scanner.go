package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/source/linemap"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

// Scanner is a token.Reader that performs lexical analysis on a stream or characters.
type Scanner struct {
	reader    *RecordingSourceReader
	linemap   *linemap.Builder
	recorder  *diagnostic.Recorder
	peeked    token.Token
	last      token.Token
	begin     source.Offset
	lineIndex source.LineIndex
}

func NewScanner(reader source.Reader) *Scanner {
	return &Scanner{
		reader:   decorateSourceReader(reader),
		linemap:  linemap.NewBuilder(),
		recorder: diagnostic.NewRecorder(),
		last:     token.NewInvalidToken("begin", token.Position{}),
		peeked:   nil,
	}
}

func NewStringScanner(input string) *Scanner {
	return NewScanner(source.NewStringReader(input))
}

func (scanner *Scanner) Pull() token.Token {
	if scanner.reader.IsExhausted() {
		return token.EndOfFile
	}
	if peeked := scanner.peeked; peeked != nil {
		scanner.peeked = nil
		return peeked
	}
	scanner.last = scanner.next()
	return scanner.last
}

func (scanner *Scanner) Peek() token.Token {
	if scanner.reader.IsExhausted() {
		return token.EndOfFile
	}
	if scanner.peeked == nil {
		scanner.peeked = scanner.next()
	}
	return scanner.peeked
}

func (scanner *Scanner) Last() token.Token {
	return scanner.last
}

func (scanner *Scanner) resetTokenRecording() {
	scanner.reader.Reset()
	scanner.begin = scanner.reader.Index()
}

func (scanner *Scanner) createInvalidToken() token.Token {
	return token.NewInvalidToken(scanner.reader.String(), scanner.currentPosition())
}

func (scanner *Scanner) incrementLineIndex() {
	scanner.reader.resetInternalIndex()
	scanner.lineIndex++
}

func (scanner *Scanner) next() token.Token {
	scanner.SkipWhitespaces()
	scanner.resetTokenRecording()
	if scanner.reader.Peek() == source.EndOfFile {
		return token.EndOfFile
	}
	return scanner.nextNonEndOfFile()
}

func (scanner *Scanner) reportError(err error) {
	scanner.recorder.Record(diagnostic.Entry{
		Kind:    &diagnostic.Error,
		Stage:   &diagnostic.LexicalAnalysis,
		Source:  scanner.reader.String(),
		Message: err.Error(),
		Position: diagnostic.Position{
			Column:    scanner.reader.internalIndex,
			LineIndex: scanner.lineIndex,
		},
	})
}

func (scanner *Scanner) nextNonEndOfFile() token.Token {
	switch next := scanner.reader.Peek(); {
	case next.IsAlphabetic():
		return scanner.ScanIdentifier()
	}
	return token.EndOfFile
}
