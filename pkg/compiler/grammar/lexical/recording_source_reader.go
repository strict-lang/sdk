package lexical

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

// recordingSourceReader decorates a input.Reader with scanning specific logic that
// helps the scanning to report more verbose errors to the diagnostics.Bag.
type recordingSourceReader struct {
	delegate input.Reader
	builder  strings.Builder
	// internalIndex is used to capture the current line-offset. The scanning resets this
	// value every time a linebreak is encountered. This value is not the same as the value
	// returned by the Index() method. It can be reset using the ResetInternalIndex() method.
	internalIndex input.Offset
}

var _ input.Reader = &recordingSourceReader{}

func decorateSourceReader(reader input.Reader) *recordingSourceReader {
	return &recordingSourceReader{
		delegate: reader,
	}
}

func (reader *recordingSourceReader) Pull() input.Char {
	next := reader.delegate.Pull()
	reader.builder.WriteRune(rune(next))
	reader.internalIndex++
	return next
}

func (reader *recordingSourceReader) Peek() input.Char {
	return reader.delegate.Peek()
}

func (reader *recordingSourceReader) Current() input.Char {
	return reader.delegate.Current()
}

func (reader *recordingSourceReader) Index() input.Offset {
	return reader.delegate.Index()
}

func (reader *recordingSourceReader) IsExhausted() bool {
	return reader.delegate.IsExhausted()
}

func (reader *recordingSourceReader) Skip(count int) {
	reader.delegate.Skip(count)
	reader.internalIndex += 2
}

func (reader *recordingSourceReader) Reset() {
	reader.builder.Reset()
}

func (reader *recordingSourceReader) resetInternalIndex() {
	reader.internalIndex = 0
}

func (reader *recordingSourceReader) String() string {
	return reader.builder.String()
}
