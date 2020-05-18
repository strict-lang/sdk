package analysis

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

type Analysis struct {
	ImportScope scope.Scope
}

func (analysis *Analysis) Store(isolate *isolate.Isolate) {
	isolate.Properties.Insert("analysis", analysis)
}

func RequireInIsolate(isolate *isolate.Isolate) *Analysis {
	if property, ok := isolate.Properties.Lookup("analysis"); ok {
		if analysis, ok := property.(*Analysis); ok {
			return analysis
		}
		log.Fatalf(`found invalid property named "analysis": %+v`, property)
	}
	isolate.Properties.Log()
	log.Fatalln("can not find scope in isolate")
	return nil
}

type Creation struct {
	Unit       *tree.TranslationUnit
	Namespaces *namespace.Table
	NamespaceSymbol *scope.Namespace
}

func (creation *Creation) Create() *Analysis {
	importScope := creation.createImportScope()
	return &Analysis{importScope}
}

func (creation *Creation) createImportScope() scope.Scope {
	namespaces := append(creation.resolveAllNamespaces(), creation.NamespaceSymbol)
	importScope := scope.NewImportScope(creation.Unit.Name, namespaces)
	namespaceScope := creation.NamespaceSymbol.Scope
	return scope.Combine(importScope.Id(), namespaceScope, importScope)
}

func (creation *Creation) listImportedNamespaces() (namespaces []string) {
	for _, statement := range creation.Unit.Imports {
		namespaces = append(namespaces, statement.Target.Namespace())
	}
	return namespaces
}

func (creation *Creation) resolveAllNamespaces() (symbols []scope.Symbol) {
	namespaces := creation.listImportedNamespaces()
	for _, symbol := range namespaces {
		symbols = append(symbols, creation.resolveNamespace(symbol))
	}
	return symbols
}

func (creation *Creation) resolveNamespace(namespace string) scope.Symbol {
	cache := scope.GlobalNamespaceTable()
	if symbol, ok := cache.Lookup(namespace); ok {
		return symbol
	}
	createdNamespace := creation.importNamespace(namespace)
	cache.Insert(namespace, createdNamespace)
	return createdNamespace
}

func (creation *Creation) importNamespace(namespace string) *scope.Namespace {
	log.Printf("don't know how to import namespace: " + namespace)
	return nil
}

