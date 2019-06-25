package source

import "bufio"

type BufIoReader struct {
	source bufio.Reader
}

func NewBufIoReader(source bufio.Reader) Reader {
	return NewBufIoReader(source)
}

func (reader *BufIoReader) Pull() Char {
	next, _, err := reader.source.ReadRune()
	if err != nil {
		return EndOfFile
	}
	return Char(next)
}

func (reader *BufIoReader) Peek() Char {
	return EndOfFile
}

