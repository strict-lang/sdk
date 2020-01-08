package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/scope"

func ResolveNearestScope(node Node) (scope.Scope, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if scopeOwner, ok := node.(ScopeOwner); ok {
			return scopeOwner.Scope(), true
		}
		currentParent, _ = currentParent.EnclosingNode()
	}
	return nil, false
}
