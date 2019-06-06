package source

type stringReader struct {
	source string
	length int
	index  int
	last   Char
	peeked Char
}

func NewStringReader(source string) Reader {
	return &stringReader{
		index:  -1, // 0 after first pull
		source: source,
		length: len(source),
		peeked: EndOfFile,
		last:   BeginOfFile,
	}
}

func (reader *stringReader) IsExhausted() bool {
	return reader.index >= reader.length
}

func (reader *stringReader) Pull() Char {
	reader.index++
	if reader.IsExhausted() {
		return EndOfFile
	}
	reader.peeked = EndOfFile
	pulled := Char(reader.source[reader.index])
	reader.last = pulled
	return pulled
}

func (reader *stringReader) Peek() Char {
	if reader.peeked != EndOfFile {
		return reader.peeked
	}
	if reader.IsExhausted() {
		return EndOfFile
	}
	if reader.index+1 >= reader.length {
		return EndOfFile
	}
	peeked := Char(reader.source[reader.index+1])
	reader.peeked = peeked
	return peeked
}

func (reader *stringReader) Last() Char {
	return reader.last
}

func (reader *stringReader) Skip(count int) {
	reader.index += count
}
