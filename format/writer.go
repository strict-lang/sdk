package format

import "strings"

type Writer interface {
	Write(string)
	WriteRune(rune)
}

type StringWriter struct {
	Buffer *strings.Builder
}

func NewStringWriter() *StringWriter {
	return &StringWriter{
		Buffer: &strings.Builder{},
	}
}

func (writer *StringWriter) String() string {
	return writer.Buffer.String()
}

func (writer *StringWriter) Write(text string) {
	writer.Buffer.WriteString(text)
}

func (writer *StringWriter) WriteRune(value rune) {
	writer.Buffer.WriteRune(value)
}