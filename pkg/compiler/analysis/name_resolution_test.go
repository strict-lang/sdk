package analysis

import (
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"strict.dev/sdk/pkg/compiler/grammar/syntax"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/isolate"
	passes "strict.dev/sdk/pkg/compiler/pass"
	"testing"
)

func parseTestUnit() *tree.TranslationUnit {
	result := syntax.ParseString("Test",
		`
method Number add(Number left, Number right)
  return left + right

method Number addPositive(Number left, Number right)
  if left < 0 or right < 0
    return 0
  return add(left, right)

`)
	return result.TranslationUnit
}

func TestNameResolutionPass(testing *testing.T) {
	context := &passes.Context{
		Unit:       parseTestUnit(),
		Diagnostic: diagnostic.NewBag(),
		Isolate:    isolate.SingleThreaded(),
	}
	execution, _ := passes.NewExecution(NameResolutionPassId, context)
	if err := execution.Run(); err != nil {
		testing.Error(err)
	}
}
