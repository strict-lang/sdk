package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/scope"

func ResolveNearestScope(node Node) (scope.Scope, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if scopeOwner, ok := currentParent.(ScopeOwner); ok && scopeOwner.Scope() != nil {
			return scopeOwner.Scope(), true
		}
		currentParent, _ = currentParent.EnclosingNode()
	}
	return nil, false
}

func ResolveNearestMutableScope(node Node) (scope.MutableScope, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if mutableScope, exists := mutableScopeOfNode(currentParent); exists {
			return mutableScope, true
		}
		currentParent, _ = currentParent.EnclosingNode()
	}
	return nil, false
}

func IsInsideOfMethod(node Node) bool {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if _, isMethod := currentParent.(*MethodDeclaration); isMethod {
			return true
		}
	}
	return false
}

func mutableScopeOfNode(node Node) (scope.MutableScope, bool) {
	if scopeOwner, ok := node.(ScopeOwner); ok {
		someScope := scopeOwner.Scope()
		if mutableScope, isMutable := someScope.(scope.MutableScope); isMutable {
			return mutableScope, true
		}
	}
	return nil, false
}
