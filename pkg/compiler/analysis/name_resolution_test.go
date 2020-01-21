package analysis

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"testing"
)

var testUnit = syntax.ParseString("Test",
	`
method Number add(Number left, Number right)
  return left + right

method Number addPositive(Number left, Number right)
  if left < 0 or right < 0 do
    return 0
  return add(left, right)

`).TranslationUnit

func TestNameResolutionPass(testing *testing.T) {
	context := &passes.Context{
		Unit:       testUnit,
		Diagnostic: diagnostic.NewBag(),
		Isolate:    isolate.SingleThreaded(),
	}
	execution, _ := passes.NewExecution(NameResolutionPassId, context)
	if err := execution.Run(); err != nil {
		testing.Error(err)
	}
	fmt.Printf("%#v\n",context.Unit)
}
