package arduino

import (
	backends "strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/backend/cpp"
)

func Generate(input backends.Input) (backends.Output, error) {
	backend := NewBackend()
	return backend.Generate(input)
}

type Backend struct {}

func NewBackend() *Backend {
	return &Backend{}
}

func (backend *Backend) Generate(input backends.Input) (backends.Output, error) {
	generation := cpp.NewGenerationWithExtension(input, NewGeneration())
	return backends.Output{
		GeneratedFiles: []backends.GeneratedFile{
			{
				Name: input.Unit.Name + ".ino",
				Content: []byte(generation.Generate()),
			},
		},
	}, nil
}
