package scanner

import (
	"gitlab.com/strict-lang/sdk/compiler/source"
	"strings"
)

// RecordingSourceReader decorates a source.Reader with scanner specific logic that
// helps the scanner to report more verbose errors to the diagnostics.Recorder.
type RecordingSourceReader struct {
	delegate source.Reader
	builder  strings.Builder
	// internalIndex is used to capture the current line-offset. The scanner resets this
	// value every time a linebreak is encountered. This value is not the same as the value
	// returned by the Index() method. It can be reset using the ResetInternalIndex() method.
	internalIndex source.Offset
}

var _ source.Reader = &RecordingSourceReader{}

func decorateSourceReader(reader source.Reader) *RecordingSourceReader {
	return &RecordingSourceReader{
		delegate: reader,
	}
}

func (reader *RecordingSourceReader) Pull() source.Char {
	next := reader.delegate.Pull()
	reader.builder.WriteRune(rune(next))
	return next
}

func (reader *RecordingSourceReader) Peek() source.Char {
	return reader.delegate.Peek()
}

func (reader *RecordingSourceReader) Last() source.Char {
	return reader.delegate.Last()
}

func (reader *RecordingSourceReader) Index() source.Offset {
	return reader.delegate.Index()
}

func (reader *RecordingSourceReader) IsExhausted() bool {
	return reader.delegate.IsExhausted()
}

func (reader *RecordingSourceReader) Skip(count int) {
	reader.delegate.Skip(count)
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
