package lowering

import (
	isolates "github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
)

const FullLoweringId passes.Id = "FullLowering"

type Lowering struct {}

func (lowering *Lowering) Dependencies(isolate *isolates.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, LetBindingLoweringPassId, ComputationCallLoweringPassId)
}

func (lowering *Lowering) Id() passes.Id {
	return FullLoweringId
}

func (lowering *Lowering) Run(context *passes.Context) {}
