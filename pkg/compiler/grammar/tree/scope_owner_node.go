package tree

import "strict.dev/sdk/pkg/compiler/scope"

type ScopeOwner interface {
	Scope() scope.Scope
	UpdateScope(target scope.Scope)
}