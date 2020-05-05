package namespace

import "github.com/strict-lang/sdk/pkg/compiler/grammar/tree"

type Namespace interface {
	QualifiedName() string
	Dependencies() []Dependency
	Entries() []Entry
}

type Dependency interface {
	Load()
}

type Entry interface {
	FileName() string
	TranslationUnit() *tree.TranslationUnit
}

type namespace struct {
	qualifiedName string
	entries  []Entry
	computed bool
}

func (namespace *namespace) QualifiedName() string {
	return namespace.qualifiedName
}

func (namespace *namespace) Entries() []Entry {
	return namespace.entries
}

func (namespace *namespace) Dependencies() []Namespace {

	return nil
}

type entry struct {
	fileName string
	translationUnit *tree.TranslationUnit
}

func (entry *entry) FileName() string {
	return entry.fileName
}

func (entry *entry) TranslationUnit() *tree.TranslationUnit {
	return entry.translationUnit
}

func NewRoot(directory string) Namespace {
	return nil
}


