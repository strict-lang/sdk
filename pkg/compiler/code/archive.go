package code

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"os"
)

type Archive interface {
	Import(scope *Scope) error
}

type singleFileArchive struct {
	path string
}

func NewSingleFileArchive(path string) Archive {
	return &singleFileArchive{path: path}
}

func (archive *singleFileArchive) Import(scope *Scope) error {
	translationUnit, err := parseFile(archive.path)
	if err != nil {
		return err
	}
	Analyse(translationUnit, scope)
	return nil
}

func parseFile(path string) (*tree.TranslationUnit, error) {
	_, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return nil, nil
}