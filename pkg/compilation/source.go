package compilation

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	"os"
)

type Source interface {
	newSourceReader() source2.Reader
}

type FileSource struct {
	File *os.File
}

func (fileSource *FileSource) newSourceReader() source2.Reader {
	return source2.NewStreamReader(fileSource.File)
}

type InMemorySource struct {
	Source string
}

func (inMemorySource *InMemorySource) newSourceReader() source2.Reader {
	return source2.NewStringReader(inMemorySource.Source)
}
