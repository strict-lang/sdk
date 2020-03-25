package compiler

import (
	"os"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type Source interface {
	newSourceReader() input.Reader
}

type FileSource struct {
	File *os.File
}

func (fileSource *FileSource) newSourceReader() input.Reader {
	return input.NewStreamReader(fileSource.File)
}

type InMemorySource struct {
	Source string
}

func (inMemorySource *InMemorySource) newSourceReader() input.Reader {
	return input.NewStringReader(inMemorySource.Source)
}
