package code

import (
	"testing"
)

type Dummy struct {
	scope *Scope
}

func (dummy Dummy) Scope() *Scope {
	return dummy.scope
}

func TestParentModifiesChild(test *testing.T) {
	scope := NewRootScope()
	child := scope.NewChild()

	err := scope.PutSymbol("a", Dummy{scope: scope})
	expectNoError(test, err)
	expectSymbol(test, scope, "a")
	expectSymbol(test, child, "a")
}

func TestChildDoesNotModifyParent(test *testing.T) {
	scope := NewRootScope()
	child := scope.NewChild()

	err := scope.PutSymbol("a", Dummy{scope: child})
	expectNoError(test, err)
	expectSymbol(test, child, "a")
	expectNoSymbol(test, scope, "a")
}

func TestChildScopeCreation(test *testing.T) {
	scope := NewRootScope()
	err := scope.PutSymbol("a", Dummy{scope: scope})
	expectNoError(test, err)

	child := scope.NewChild()
	expectName(test, child, "@child-1")
	expectSymbol(test, child, "a")

	namedChild := scope.NewNamedChild("test")
	expectName(test, namedChild, "@test")
	expectSymbol(test, namedChild, "a")
}

func expectName(test *testing.T, scope *Scope, name string) {
	if scope.name == name {
		return
	}
	test.Errorf("expected name to be %s but got %s", name, scope.name)
}

func expectSymbol(test *testing.T, scope *Scope, key string) {
	if scope.ContainsSymbol(key) {
		return
	}
	test.Errorf("expected to find %s in scope %s", key, scope.name)
}

func expectNoSymbol(test *testing.T, scope *Scope, key string) {
	if !scope.ContainsSymbol(key) {
		return
	}
	test.Errorf("expected not to find %s in scope %s", key, scope.name)
}

func expectNoError(test *testing.T, err error) {
	if err != nil {
		test.Errorf("unexpected error: %s", err)
	}
}
