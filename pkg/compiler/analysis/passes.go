package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	"gitlab.com/strict-lang/sdk/pkg/compiler/pass"
)

func registerPassInstance(instance pass.Pass) {
	isolate.RegisterConfigurator(func(isolate *isolate.Isolate) {
		name := string(instance.Id())
		isolate.Properties.Insert(name, instance)
	})
}
