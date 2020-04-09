package silk

import backends "github.com/strict-lang/sdk/pkg/compiler/backend"

type Backend struct {
}

func (backend *Backend) Compile(input backends.Input) backends.Output {
	return backends.Output{}
}
