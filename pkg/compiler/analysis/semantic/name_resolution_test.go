package semantic

import (
	"github.com/strict-lang/sdk/pkg/compiler/analysis"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	isolates "github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
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

method testing(any Any)
  let hashCode = any.CalculateHashCode()

`)
	if result.Error != nil {
		log.Printf("failed to parse Unit: %v\n", result.Error)
	}
	return result.TranslationUnit
}

func TestNameResolutionPass(testing *testing.T) {
	isolate := isolates.New()
	testAnalysis := analysis.Analysis{
		ImportScope: createImportScope(),
	}
	testAnalysis.Store(isolate)
	context := &passes.Context{
		Unit:       parseTestUnit(),
		Diagnostic: diagnostic.NewBag(),
		Isolate: isolate,
	}
	execution, _ := passes.NewExecution(NameResolutionPassId, context)
	if err := execution.Run(); err != nil {
		testing.Error(err)
	}
	pretty.Print(context.Unit)
}

func createImportScope() scope.Scope {
	testScope := scope.NewOuterScope(scope.Id("test-scope"), scope.NewBuiltinScope())
	testScope.Insert(&scope.Class{
		DeclarationName: "Test",
		QualifiedName:   "Test",
	})
	return testScope
}
