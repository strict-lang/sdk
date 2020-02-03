package backend

import (
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
)

type GeneratedFile struct {
	Name string
	Content []byte
}

type Output struct {
	GeneratedFiles[] GeneratedFile
}

type Input struct {
	Unit *tree.TranslationUnit
	Diagnostics *diagnostic.Bag
}

type Backend interface {
	Generate(Input) (Output, error)
}