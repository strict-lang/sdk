package lexical

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

// RecordingSourceReader decorates a input.Reader with scanning specific logic that
// helps the scanning to report more verbose errors to the diagnostics.Bag.
type RecordingSourceReader struct {
	delegate input.Reader
	builder  strings.Builder
	// internalIndex is used to capture the current line-offset. The scanning resets this
	// value every time a linebreak is encountered. This value is not the same as the value
	// returned by the Index() method. It can be reset using the ResetInternalIndex() method.
	internalIndex input.Offset
}

var _ input.Reader = &RecordingSourceReader{}

func decorateSourceReader(reader input.Reader) *RecordingSourceReader {
	return &RecordingSourceReader{
		delegate: reader,
	}
}

func (reader *RecordingSourceReader) Pull() input.Char {
	next := reader.delegate.Pull()
	reader.builder.WriteRune(rune(next))
	reader.internalIndex++
	return next
}

func (reader *RecordingSourceReader) Peek() input.Char {
	return reader.delegate.Peek()
}

func (reader *RecordingSourceReader) Current() input.Char {
	return reader.delegate.Current()
}

func (reader *RecordingSourceReader) Index() input.Offset {
	return reader.delegate.Index()
}

func (reader *RecordingSourceReader) IsExhausted() bool {
	return reader.delegate.IsExhausted()
}

func (reader *RecordingSourceReader) Skip(count int) {
	reader.delegate.Skip(count)
	reader.internalIndex += 2
}

func (reader *RecordingSourceReader) Reset() {
	reader.builder.Reset()
}

func (reader *RecordingSourceReader) resetInternalIndex() {
	reader.internalIndex = 0
}

func (reader *RecordingSourceReader) String() string {
	return reader.builder.String()
}
