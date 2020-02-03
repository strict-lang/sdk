package analysis

import (
	"strict.dev/sdk/pkg/compiler/pass"
)

func registerPassInstance(instance pass.Pass) {
	pass.Register(instance)
}
