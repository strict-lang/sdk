package ilasm

import backends "strict.dev/sdk/pkg/compiler/backend"

type Backend struct {}

func (backend *Backend) Generate(input backends.Input) backends.Output {
	return backends.Output{
		GeneratedFiles: []backends.GeneratedFile{},
	}
}

type Generation struct {
	code *CodeBuilder
	currentClass *Class
}
