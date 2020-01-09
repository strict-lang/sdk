package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	"gitlab.com/strict-lang/sdk/pkg/compiler/pass"
)

func registerPassInstance(id string, instance pass.Pass) {
	isolate.RegisterConfigurator(func(isolate *isolate.Isolate) {
		isolate.Properties.Insert(id, instance)
	})
}
