package analysis

import (
	"github.com/strict-lang/sdk/pkg/compiler/pass"
)


func RunEntering(context *pass.Context) error {
	return pass.RunWithId(NameResolutionPassId, context)
}

func Run(context *pass.Context) error {
	return pass.RunWithId(NameResolutionPassId, context)
}


func registerPassInstance(instance pass.Pass) {
	pass.Register(instance)
}
