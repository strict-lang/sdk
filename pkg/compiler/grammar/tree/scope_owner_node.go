package tree

import "github.com/strict-lang/sdk/pkg/compiler/scope"

type ScopeOwner interface {
	Scope() scope.Scope
	UpdateScope(target scope.Scope)
}
