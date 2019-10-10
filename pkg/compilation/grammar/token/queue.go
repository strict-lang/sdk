package token

type Queue []Token

type QueueReader struct {
	index int
	queue Queue
}

var _ Stream = &QueueReader{}

func NewQueueReader(queue Queue) *QueueReader {
	return &QueueReader{
		index: -1,
		queue: queue,
	}
}

func (reader *QueueReader) hasNext() bool {
	return len(reader.queue) > reader.index+1
}

func (reader *QueueReader) Pull() Token {
	if !reader.hasNext() {
		return EndOfFile
	}
	reader.index++
	element := reader.queue[reader.index]
	return element
}

func (reader *QueueReader) Peek() Token {
	if !reader.hasNext() {
		return EndOfFile
	}
	return reader.queue[reader.index+1]
}

func (reader *QueueReader) Last() Token {
	if reader.index < 0 {
		return EndOfFile
	}
	if !reader.hasNext() {
		return EndOfFile
	}
	return reader.queue[reader.index]
}
