package compiler

import (
	"gitlab.com/strict-lang/sdk/compiler/source"
	"os"
)

type Source interface {
	newSourceReader() source.Reader
}

type FileSource struct {
	File *os.File
}

func (fileSource *FileSource) newSourceReader() source.Reader {
	return source.NewStreamReader(fileSource.File)
}

type InMemorySource struct {
	Source string
}

func (inMemorySource *InMemorySource) newSourceReader() source.Reader {
	return source.NewStringReader(inMemorySource.Source)
}
