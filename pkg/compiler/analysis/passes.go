package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/pass"
)

func registerPassInstance(instance pass.Pass) {
	pass.Register(instance)
}
