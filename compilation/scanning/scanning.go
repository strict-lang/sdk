package scanning

import (
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/source/linemap"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

const (
	TabIndent        token.Indent = 4
	WhitespaceIndent token.Indent = 1
)

// Scanning is a token.Stream that performs lexical analysis on a stream or characters.
type Scanning struct {
	input          *RecordingSourceReader
	lineMapBuilder *linemap.Builder
	diagnosticBag  *diagnostic.Bag
	// peeked points to the most recently peeked token.
	peeked token.Token
	// last points to the most recently scanned token. It is an InvalidToken if no other
	// token has been scanned. The fields value is never nil.
	last token.Token
	// begin is the begin index of the token that is currently scanned. It is set to the
	// current offset when the scanning starts scanning the next token.
	begin source.Offset
	// lineIndex current lineIndex of the scanning, incremented each time a linefeed is hit.
	// The scanning keeps track of his line-index to report better errors to the diagnostics.
	lineIndex source.LineIndex
	// endOfStatementPrevention tells the scanning whether it should insert an EndOfStatement or not.
	// It will only insert an EndOfStatement token when this field is zero. This is not a boolean,
	// because preventions can be nested. Common tokens that prevent the scanning from generating
	// an EndOfStatement token are: Parentheses and Brackets
	endOfStatementPrevention int
	// indent is the current indentation level. It is updates while scanning and assigned
	// to all tokens that are created.
	indent token.Indent
	// updateIndent is a flag that tells the scanning whether it should update the indent
	// value. It is set and unset during scanning. Once the first non-whitespace character
	// in a line is hit, the scanning disables this flag. All scanned tokens in that line
	// will get the indent value of the 'ident' field, which can not change anymore. Once
	// a linefeed is hit, the indent is reset and the updateIndent is set to true.
	updateIndent bool
	// emptyLine records whether the currently scanned line is empty. If it is, the scanning
	// will not insert an EndOfStatement token even if 'insertEos' is set to true.
	emptyLine bool
	// lineBeginOffset is the offset to the lines begin. It is updated each
	// time a new line is added to the lineMapBuilder.
	lineBeginOffset source.Offset
}

func NewDiagnosticScanner(reader source.Reader, recorder *diagnostic.Bag) *Scanning {
	beginOfFile := token.NewInvalidToken("BeginOfFile", token.Position{}, token.NoIndent)
	scanning := &Scanning{
		input:          decorateSourceReader(reader),
		lineMapBuilder: linemap.NewBuilder(),
		diagnosticBag:  recorder,
		last:           beginOfFile,
		lineIndex:      1,
		peeked:         nil,
		updateIndent:   true,
		emptyLine:      true, // The line is empty until a char is hit
	}
	scanning.advance()
	return scanning
}

func NewScanning(reader source.Reader) *Scanning {
	return NewDiagnosticScanner(reader, diagnostic.NewBag())
}

var _ token.Stream = &Scanning{}

func NewStringScanning(input string) *Scanning {
	return NewScanning(source.NewStringReader(input))
}

func (scanning *Scanning) advance() {
	scanning.input.Pull()
}

func (scanning *Scanning) char() source.Char {
	return scanning.input.Current()
}

func (scanning *Scanning) peekChar() source.Char {
	return scanning.input.Peek()
}

func (scanning *Scanning) Pull() token.Token {
	if scanning.input.IsExhausted() {
		return scanning.endOfFile()
	}
	if peeked := scanning.peeked; peeked != nil {
		scanning.last = peeked
		scanning.peeked = nil
		return peeked
	}
	next := scanning.next()
	scanning.last = next
	return next
}

func (scanning *Scanning) Peek() token.Token {
	if scanning.peeked == nil {
		next := scanning.next()
		scanning.peeked = next
		return next
	}
	return scanning.peeked
}

// endOfFile returns either an EndOfStatement or an EndOfFile token.
// If there was no final-newline and therefor no final end-of-statement, the scanning
// will first return an end-of-statement. There will never be two end-of-statements
// at the end of a file.
func (scanning *Scanning) endOfFile() token.Token {
	last := scanning.last
	if _, ok := last.(*token.EndOfStatementToken); ok {
		return token.EndOfFile
	}
	newLast := token.NewEndOfStatementToken(scanning.offset())
	scanning.last = newLast
	return scanning.last
}

func (scanning *Scanning) Last() token.Token {
	return scanning.last
}

func (scanning *Scanning) resetTokenRecording() {
	scanning.input.Reset()
	scanning.begin = scanning.input.Index()
}

func (scanning *Scanning) createInvalidToken() token.Token {
	return token.NewInvalidToken(
		scanning.input.String(), scanning.currentPosition(), scanning.indent)
}

func (scanning *Scanning) incrementLineIndex() (token.Token, bool) {
	scanning.indent = 0
	scanning.updateIndent = true
	length := scanning.offset() - scanning.lineBeginOffset
	scanning.lineMapBuilder.Append(scanning.lineBeginOffset, length)
	scanning.input.resetInternalIndex()
	scanning.lineIndex++
	scanning.lineBeginOffset = scanning.offset()
	if !scanning.shouldInsertEndOfStatement() || scanning.emptyLine {
		return nil, false
	}
	scanning.emptyLine = true
	return token.NewEndOfStatementToken(scanning.offset()), true
}

func (scanning *Scanning) shouldInsertEndOfStatement() bool {
	return scanning.endOfStatementPrevention == 0
}

func (scanning *Scanning) next() token.Token {
	if endOfStatement, ok := scanning.skipWhitespaces(); ok {
		// The SkipWhitespaces method returns an EndOfStatementToken if it hits a
		// linefeed character while the scanners 'insertEos' flag is set.
		return endOfStatement
	}
	scanning.resetTokenRecording()
	scanning.updateIndent = false
	scanning.emptyLine = false
	if scanning.char() == source.EndOfFile {
		return scanning.endOfFile()
	}
	return scanning.nextNonEndOfFile()
}

func (scanning *Scanning) reportError(err error) {
	scanning.diagnosticBag.Record(diagnostic.RecordedEntry{
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.LexicalAnalysis,
		Message:  err.Error(),
		Position: scanning.last.Position(),
	})
}

func (scanning *Scanning) SkipComment() {
	// Skip the next '/' characters
	scanning.advance()
	scanning.advance()
	for {
		scanning.advance()
		if scanning.char() == '\n' {
			break
		}
	}
}

func (scanning *Scanning) nextNonEndOfFile() token.Token {
	switch next := scanning.char(); {
	case next == '/':
		if scanning.peekChar() == '/' {
			scanning.SkipComment()
			return scanning.nextNonEndOfFile()
		}
	case next.IsAlphabetic():
		return scanning.ScanIdentifierOrKeyword()
	case next.IsNumeric():
		return scanning.ScanNumber()
	case isKnownOperator(next):
		return scanning.ScanOperator()
	case next == '"':
		return scanning.ScanStringLiteral()
	}
	if scanning.input.IsExhausted() {
		return scanning.endOfFile()
	}
	return scanning.createInvalidToken()
}

func (scanning *Scanning) NewLineMap() *linemap.LineMap {
	return scanning.lineMapBuilder.NewLineMap()
}
