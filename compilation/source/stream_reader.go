package source

import (
	"bufio"
	"io"
)

type streamReader struct {
	index     Offset
	stream    *bufio.Reader
	current   Char
	peeked    Char
	exhausted bool
}

func (reader *streamReader) Pull() Char {
	if reader.exhausted {
		return EndOfFile
	}
	reader.peeked = EndOfFile
	reader.index++
	raw, _, err := reader.stream.ReadRune()
	if err != nil {
		reader.exhausted = true
		return EndOfFile
	}
	casted := Char(raw)
	reader.current = casted
	return casted
}

func (reader *streamReader) Peek() Char {
	char, _, err := reader.stream.ReadRune()
	if err != nil {
		return EndOfFile
	}
	if err := reader.stream.UnreadRune(); err != nil {
		return EndOfFile
	}
	reader.peeked = Char(char)
	return reader.peeked
}

func (reader *streamReader) Skip(count int) {
	reader.index += Offset(count)
	reader.stream.Discard(count)
}

func (reader *streamReader) Current() Char {
	return reader.current
}

func (reader *streamReader) Index() Offset {
	return reader.index
}

func (reader *streamReader) IsExhausted() bool {
	return reader.exhausted
}

func NewStreamReader(reader io.Reader) Reader {
	stream := bufio.NewReader(reader)
	return &streamReader{
		stream:  stream,
		current: EndOfFile,
		peeked:  EndOfFile,
	}
}
