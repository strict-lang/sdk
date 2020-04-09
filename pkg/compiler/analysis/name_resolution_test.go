package analysis

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"testing"
)

func parseTestUnit() *tree.TranslationUnit {
	result := syntax.ParseString("Test",
		`
method add(left Number, right Number) returns Number
  return left + right

method addPositive(left Number, right Number) returns Number
  if left < 0 or right < 0
    return 0
  return add(left, right)

`)
	if result.Error != nil {
		_ = fmt.Errorf("failed to parse unit: %v\n", result.Error)
	}
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
