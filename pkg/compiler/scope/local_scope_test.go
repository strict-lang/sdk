package scope

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestLocalScope_Id(testing *testing.T) {
	scope := createTestLocalScope("test")
	if scope.Id() != "test" {
		testing.Errorf("invalid id: %s", scope.Id())
	}
}

func TestLocalScope_Insert(testing *testing.T) {
	scope := createTestLocalScope("test")
	scope.Insert(NewPositionedTestSymbol("a", 0))
	scope.Insert(NewPositionedTestSymbol("a", 1))
	point := NewReferencePoint("a")
	if symbols := scope.Lookup(point); !symbols.IsEmpty() {
		if len(symbols) == 2 {
			testing.Error("Scope has duplicated entries")
		}
		if symbol := symbols.First(); symbol.position != 1 {
			testing.Error("Scope was not overridden")
		}
	}
}

func createTestLocalScope(id Id) MutableScope {
	return NewShadowingLocalScope(
		id,
		input.CreateEmptyRegion(0),
		NewEmptyScope(""))
}
