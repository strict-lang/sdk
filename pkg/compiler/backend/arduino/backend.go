package arduino

import (
	backends "github.com/strict-lang/sdk/pkg/compiler/backend"
	"github.com/strict-lang/sdk/pkg/compiler/backend/cpp"
)

const BackendName = "arduino"

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
	generation := cpp.NewGenerationWithExtension(input, NewGeneration())
	return backends.Output{
		GeneratedFiles: []backends.GeneratedFile{
			{
				Name:    input.Unit.Name + ".ino",
				Content: []byte(generation.Generate()),
			},
		},
	}, nil
}
