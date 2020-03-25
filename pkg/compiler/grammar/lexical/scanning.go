package lexical

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"strings"
)

const (
	tabIndent        token.Indent = 4
	whitespaceIndent token.Indent = 1
)

// Scanning is a token.Stream that performs lexical analysis on a stream or characters.
type Scanning struct {
	input          *recordingSourceReader
	lineMapBuilder *linemap.Builder
	diagnosticBag  *diagnostic.Bag
	// peeked points to the most recently peeked token.
	peeked token.Token
	// last points to the most recently scanned token. It is an InvalidToken if no other
	// token has been scanned. The fields value is never nil.
	last token.Token
	// begin is the begin index of the token that is currently scanned. It is set to the
	// current offset when the scanning starts scanning the next token.
	begin input.Offset
	// lineIndex current lineIndex of the scanning, incremented each time a linefeed is hit.
	// The scanning keeps track of his line-index to report better errors to the diagnostics.
	lineIndex input.LineIndex
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
	lineBeginOffset input.Offset
	hasHitEndOfFile bool
	lineBuffer      *strings.Builder
}

var beginOfFile = token.NewInvalidToken("BeginOfFile", token.Position{}, token.NoIndent)

// newDiagnosticScanning creates a Scanning that reads from the reader and
// reports errors to the given DiagnosticBag.
func newDiagnosticScanning(reader input.Reader, recorder *diagnostic.Bag) *Scanning {
	scanning := &Scanning{
		input:          decorateSourceReader(reader),
		lineMapBuilder: linemap.NewBuilder(),
		diagnosticBag:  recorder,
		last:           beginOfFile,
		lineIndex:      1,
		updateIndent:   true,
		emptyLine:      true, // The line is empty until a char is hit
		lineBuffer:     &strings.Builder{},
	}
	scanning.advance()
	return scanning
}

// NewScanning creates a Scanning that reads from the given reader.
func NewScanning(reader input.Reader) *Scanning {
	return newDiagnosticScanning(reader, diagnostic.NewBag())
}

// NewStringScanning creates a Scanning that reads from the string.
func NewStringScanning(source string) *Scanning {
	return NewScanning(input.NewStringReader(source))
}

// Peek returns the next token without advancing in position. If no next token
// was peeked, it scans one. Subsequent calls to this method without a call to
// Pull will yield the same token.
func (scanning *Scanning) Peek() token.Token {
	if scanning.peeked == nil {
		next := scanning.next()
		scanning.peeked = next
		return next
	}
	return scanning.peeked
}

// Pull returns the current token and advances to the next one. If no token was
// peeked, it scans one. An EOF token is returned, if the input is exhausted.
func (scanning *Scanning) Pull() token.Token {
	if peeked := scanning.peeked; peeked != nil {
		scanning.last = peeked
		scanning.peeked = nil
		return peeked
	}
	return scanning.scanNextToken()
}

func (scanning *Scanning) advance() {
	scanning.maybeWriteCurrentCharacter()
	scanning.input.Pull()
}

func (scanning *Scanning) maybeWriteCurrentCharacter() {
	current := scanning.char()
	if !current.IsLineFeed() {
		scanning.lineBuffer.WriteRune(rune(current))
	}
}

// char returns the character at the inputs current position.
func (scanning *Scanning) char() input.Char {
	return scanning.input.Current()
}

// peekChar returns the next character in the input.
func (scanning *Scanning) peekChar() input.Char {
	return scanning.input.Peek()
}

func (scanning *Scanning) scanNextToken() token.Token {
	if scanning.input.IsExhausted() {
		return scanning.createEndOfFile()
	}
	next := scanning.next()
	scanning.last = next
	return next
}

// createEndOfFile returns either an EndOfStatement or an EndOfFile token.
// If there was no final-newline and therefor no final end-of-statement, the scanning
// will first return an end-of-statement. There will never be two end-of-statements
// at the end of a file.
func (scanning *Scanning) createEndOfFile() token.Token {
	if scanning.hasHitEndOfFile {
		return token.EndOfFile
	}
	if token.IsEndOfStatementToken(scanning.last) {
		scanning.hasHitEndOfFile = true
		scanning.last = token.EndOfFile
	} else {
		newLast := token.NewEndOfStatementToken(scanning.offset())
		scanning.last = newLast
	}
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

func (scanning *Scanning) advanceLine() (token.Token, bool) {
	scanning.saveCurrentLine()
	scanning.resetLineStats()
	scanning.lineIndex++
	scanning.lineBeginOffset = scanning.offset()
	return scanning.completeCurrentLine()
}

func (scanning *Scanning) completeCurrentLine() (token.Token, bool) {
	if scanning.emptyLine || !scanning.shouldInsertEndOfStatement() {
		return nil, false
	}
	scanning.emptyLine = true
	return token.NewEndOfStatementToken(scanning.offset()), true
}

// saveCurrentLine saves the characters of the current line to the linemap.
func (scanning *Scanning) saveCurrentLine() {
	length := scanning.offset() - scanning.lineBeginOffset
	text := scanning.lineBuffer.String()
	scanning.lineMapBuilder.Append(text, scanning.lineBeginOffset, length)
}

// resetLineStats clears the information of the last line.
func (scanning *Scanning) resetLineStats() {
	scanning.indent = 0
	scanning.updateIndent = true
	scanning.input.resetInternalIndex()
	scanning.lineBuffer.Reset()
}

// shouldInsertEndOfStatement determines whether an implicit EOS should be
// inserted. This is true when there are no open parentheses.
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
	if scanning.char() == input.EndOfFile || scanning.input.IsExhausted() {
		return scanning.createEndOfFile()
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

func (scanning *Scanning) skipComment() {
	scanning.skipCommentPrefix()
	scanning.skipCommentContent()
}

const commentPrefixLength = 2

func (scanning *Scanning) skipCommentPrefix() {
	scanning.tryToSkipMultiple('/', commentPrefixLength)
}

func (scanning *Scanning) skipCommentContent() {
	for {
		scanning.advance()
		if scanning.char() == '\n' {
			scanning.advance()
			scanning.skipWhitespaces()
			scanning.resetTokenRecording()
			scanning.advanceLine()
			break
		}
	}
}

func (scanning *Scanning) isLookingAtComment() bool {
	return scanning.char() == '/' && scanning.peekChar() == '/'
}

func (scanning *Scanning) nextNonEndOfFile() token.Token {
	if scanning.isLookingAtComment() {
		scanning.skipComment()
		return scanning.nextNonEndOfFile()
	}
	return scanning.scanToken()
}

func (scanning *Scanning) scanToken() token.Token {
	switch next := scanning.char(); {
	case next == '\n' || next == '\r':
		scanning.advance()
		return scanning.nextNonEndOfFile()
	case next.IsAlphabetic():
		return scanning.scanIdentifierOrKeyword()
	case next.IsNumeric():
		return scanning.scanNumber()
	case isKnownOperator(next):
		return scanning.scanOperator()
	case next == '"':
		return scanning.scanStringLiteral()
	}
	if scanning.input.IsExhausted() {
		return scanning.createEndOfFile()
	}
	return scanning.createInvalidToken()
}

// NewLineMap completes the LineMap that has been built by the scanning.
func (scanning *Scanning) NewLineMap() *linemap.LineMap {
	return scanning.lineMapBuilder.NewLineMap()
}

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
				return scanning.createEndOfFile(), true
			}
			if endOfStatement, ok := scanning.advanceLine(); ok {
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
	scanning.addIndent(tabIndent)
}

func (scanning *Scanning) addWhitespaceIndent() {
	scanning.addIndent(whitespaceIndent)
}

func (scanning *Scanning) addIndent(indent token.Indent) {
	if scanning.updateIndent {
		scanning.indent += indent
	}
}

func scanRemaining(scanning *Scanning) (tokens []token.Token) {
	for {
		next := scanning.Pull()
		if token.IsEndOfFileToken(next) {
			break
		}
		tokens = append(tokens, next)
	}
	return
}

type unexpectedCharError struct {
	got      input.Char
	expected string
}

func (err unexpectedCharError) Error() string {
	return fmt.Sprintf("unexpected char '%c', expected '%s'", err.got, err.expected)
}
