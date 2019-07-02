package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/diagnostic"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/source/linemap"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

const (
	TabIndent 			 token.Indent = 2
	WhitespaceIndent token.Indent= 1
)

// Scanner is a token.Reader that performs lexical analysis on a stream or characters.
type Scanner struct {
	reader    *RecordingSourceReader
	linemap   *linemap.Builder
	recorder  *diagnostic.Recorder
	// peeked points to the most recently peeked token.
	peeked    token.Token
	// last points to the most recently scanned token. It is an InvalidToken if no other
	// token has been scanned. The fields value is never nil.
	last      token.Token
	// begin is the begin index of the token that is currently scanned. It is set to the
	// current offset when the scanner starts scanning the next token.
	begin     source.Offset
	// lineIndex current lineIndex of the scanner, incremented each time a linefeed is hit.
	// The scanner keeps track of his line-index to report better errors to the diagnostics.
	lineIndex source.LineIndex
	// endOfStatementPrevention tells the scanner whether it should insert an EndOfStatement or not.
	// It will only insert an EndOfStatement token when this field is zero. This is not a boolean,
	// because preventers can be nested. Common tokens that prevent the scanner from generating
	// an EndOfStatement token are: Parentheses and Brackets
	endOfStatementPrevention int
	// indent is the current indentation level. It is updates while scanning and assigned
	// to all tokens that are created.
	indent   token.Indent
	// updateIndent is a flag that tells the scanner whether it should update the indent
	// value. It is set and unset during scanning. Once the first non-whitespace character
	// in a line is hit, the scanner disables this flag. All scanned tokens in that line
	// will get the indent value of the 'ident' field, which can not change anymore. Once
	// a linefeed is hit, the indent is reset and the updateIndent is set to true.
	updateIndent bool
	// emptyLine records whether the currently scanned line is empty. If it is, the scanner
	// will not insert an EndOfStatement token even if 'insertEos' is set to true.
	emptyLine bool
}

func NewScanner(reader source.Reader) *Scanner {
	return &Scanner{
		reader:   decorateSourceReader(reader),
		linemap:  linemap.NewBuilder(),
		recorder: diagnostic.NewRecorder(),
		last:     token.NewAnonymousInvalidToken(),
		peeked:   nil,
		updateIndent: true,
		emptyLine: true, // The line is empty until a char is hit
	}
}

func NewStringScanner(input string) *Scanner {
	return NewScanner(source.NewStringReader(input))
}

func (scanner *Scanner) Pull() token.Token {
	if scanner.reader.IsExhausted() {
		return scanner.endOfFile()
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
		return scanner.endOfFile()
	}
	if scanner.peeked == nil {
		scanner.peeked = scanner.next()
	}
	return scanner.peeked
}

// endOfFile returns either an EndOfStatement or an EndOfFile token.
// If there was no final-newline and therefor no final end-of-statement, the scanner
// will first return an end-of-statement. There will never be two end-of-statements
// at the end of a file.
func (scanner *Scanner) endOfFile() token.Token {
	if _, ok := scanner.last.(*token.EndOfStatementToken); ok {
		return token.EndOfFile
	}
	last := token.NewEndOfStatementToken(scanner.offset())
	scanner.last = last
	return last
}

func (scanner *Scanner) Last() token.Token {
	return scanner.last
}

func (scanner *Scanner) resetTokenRecording() {
	scanner.reader.Reset()
	scanner.begin = scanner.reader.Index()
}

func (scanner *Scanner) createInvalidToken() token.Token {
	return token.NewInvalidToken(scanner.reader.String(), scanner.currentPosition(), scanner.indent)
}

func (scanner *Scanner) incrementLineIndex() (token.Token, bool) {
	scanner.indent = 0
	scanner.updateIndent = true
	scanner.reader.resetInternalIndex()
	scanner.lineIndex++
	if !scanner.shouldInsertEndOfStatement() || scanner.emptyLine {
		return nil, false
	}
	scanner.emptyLine = true
	return token.NewEndOfStatementToken(scanner.offset()), true
}

func (scanner *Scanner) shouldInsertEndOfStatement() bool {
	return scanner.endOfStatementPrevention == 0
}

func (scanner *Scanner) next() token.Token {
	if endOfStatement, ok := scanner.SkipWhitespaces(); ok {
		// The SkipWhitespaces method returns an EndOfStatementToken if it hits a
		// linefeed character while the scanners 'insertEos' flag is set.
		return endOfStatement
	}
	scanner.resetTokenRecording()
	scanner.updateIndent = false
	scanner.emptyLine = false
	if scanner.reader.Peek() == source.EndOfFile {
		return scanner.endOfFile()
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
	case next.IsNumeric():
		return scanner.ScanNumber()
	case isKnownOperator(next):
		return scanner.ScanOperator()
	case next == '"':
		return scanner.ScanStringLiteral()
	}
	return scanner.createInvalidToken()
}