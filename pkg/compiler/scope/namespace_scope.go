package scope

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
)

func NewNamespaceScope(namespace namespace.Namespace, classes []*Class) Scope {
	id := Id("namespace." + namespace.QualifiedName())
	scope := NewOuterScopeWithRootId(id, builtinScope)
	creation := &namespaceScopeCreation{
		scope: scope,
		classes: classes,
	}
	creation.run()
	return scope
}

type namespaceScopeCreation struct {
	classes []*Class
	scope MutableScope
}

func (creation *namespaceScopeCreation) run() {
	creation.populateContents()
}

func (creation *namespaceScopeCreation) populateContents() {
	for _, class := range creation.classes {
		creation.scope.Insert(class)
	}
}
