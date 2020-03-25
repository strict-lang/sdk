package ilasm

import backends "strict.dev/sdk/pkg/compiler/backend"

const BackendName = "msil"

func init() {
	backends.Register(BackendName, NewBackend)
}

type Backend struct{}

func NewBackend() backends.Backend {
	return &Backend{}
}

func (backend *Backend) Generate(input backends.Input) (backends.Output, error) {
	return backends.Output{
		GeneratedFiles: []backends.GeneratedFile{},
	}, nil
}

type Generation struct {
	code              *BlockBuilder
	currentClass      *Class
	method            *MethodContext
	breakLabel        *Label
	continuationLabel *Label
}

func (generation *Generation) updateCurrentBlock(target *BlockBuilder) {}
