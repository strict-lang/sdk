package scope

import "log"

func NewImportScope(name string, symbols []Symbol) Scope {
	id := Id("import." + name)
	scope := NewOuterScopeWithRootId(id, builtinScope)
	creation := &importScopeCreation{
		symbols:  symbols,
		scope: scope,
	}
	creation.run()
	return scope
}

type importScopeCreation struct {
	symbols []Symbol
	scope MutableScope
}

func (creation *importScopeCreation) run() {
	creation.populateContents()
}

func (creation *importScopeCreation) populateContents() {
	for _, symbol := range creation.symbols {
		creation.insert(symbol)
	}
}

func (creation *importScopeCreation) insert(symbol Symbol) {
	if namespace, ok := AsNamespaceSymbol(symbol); ok {
		creation.insertNamespace(namespace)
		return
	}
	log.Printf("ignoring symbol %v: it is not of type namespace", symbol)
}

func (creation *importScopeCreation) insertNamespace(namespace *Namespace) {
	creation.scope.Insert(namespace)
	creation.findAndInsertTopClass(namespace)
}

func (creation *importScopeCreation) findAndInsertTopClass(namespace *Namespace) {
	classes := namespace.Scope.Search(filterForClass)
	for _, symbol := range classes {
		if class, ok := AsClassSymbol(symbol.Symbol); ok && isTopClass(namespace, class) {
			creation.insertTopClass(namespace, class)
			break
		}
	}
}

func (creation *importScopeCreation) insertTopClass(namespace *Namespace, class *Class) {
	creation.scope.Insert(class)
}

func filterForClass(symbol Symbol) bool {
	_, isClass := symbol.(*Class)
	return isClass
}

