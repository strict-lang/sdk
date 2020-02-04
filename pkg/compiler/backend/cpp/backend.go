package cpp

import (
	backends "strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
)

const BackendName = "c++"

func init() {
	backends.Register(BackendName, NewBackend)
}

func Generate(input backends.Input) (backends.Output, error) {
	backend := NewBackend()
	return backend.Generate(input)
}

type Backend struct{}

func NewBackend() backends.Backend {
	return &Backend{}
}

func (backend *Backend) Generate(input backends.Input) (backends.Output, error) {
	files := make(chan backends.GeneratedFile)
	go func() {
		files <- backend.generateHeaderFile(input)
	}()
	go func() {
		files <- backend.generateSourceFile(input)
	}()
	var output []backends.GeneratedFile
	if isContainingTestDefinitions(input.Unit) {
		output = append(output, backend.generateTestFile(input))
	}
	output = append(output, <-files)
	output = append(output, <-files)
	return backends.Output{GeneratedFiles: output}, nil
}

func isContainingTestDefinitions(node tree.Node) bool {
	counter := tree.NewCounter()
	counter.Count(node)
	return counter.ValueFor(tree.TestStatementNodeKind) != 0
}

func (backend *Backend) generateHeaderFile(input backends.Input) backends.GeneratedFile {
	generation := NewGenerationWithExtension(input, NewHeaderFileGeneration())
	return backends.GeneratedFile{
		Name:    input.Unit.Name + ".h",
		Content: []byte(generation.Generate()),
	}
}

func (backend *Backend) generateSourceFile(input backends.Input) backends.GeneratedFile {
	generation := NewGenerationWithExtension(input, NewSourceFileGeneration())
	return backends.GeneratedFile{
		Name:    input.Unit.Name + ".cc",
		Content: []byte(generation.Generate()),
	}
}

func (backend *Backend) generateTestFile(input backends.Input) backends.GeneratedFile {
	generation := NewGenerationWithExtension(input, NewTestFileGeneration())
	return backends.GeneratedFile{
		Name:    input.Unit.Name + "_test.cc",
		Content: []byte(generation.Generate()),
	}
}
