package scope

func NewNamespaceScope(namespace *Namespace) Scope {
	id := Id("namespace." + namespace.QualifiedName)
	scope := NewOuterScopeWithRootId(id, builtinScope)
	creation := &namespaceScopeCreation{
		namespace: namespace,
		scope:     scope,
	}
	creation.run()
	return scope
}

type namespaceScopeCreation struct {
	namespace *Namespace
	scope MutableScope
}

func (creation *namespaceScopeCreation) run() {
	creation.populateContents()
}

func (creation *namespaceScopeCreation) populateContents() {
	for _, child := range creation.scope.Search(findAll) {
		creation.scope.Insert(child.Symbol)
	}
}

func findAll(symbol Symbol) bool {
	return true
}
