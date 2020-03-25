package backend

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
)

type GeneratedFile struct {
	Name    string
	Content []byte
}

type Output struct {
	GeneratedFiles [] GeneratedFile
}

type Input struct {
	Unit        *tree.TranslationUnit
	Diagnostics *diagnostic.Bag
}

type Backend interface {
	Generate(Input) (Output, error)
}

func LookupInIsolate(isolate *isolate.Isolate, name string) (Backend, bool) {
	propertyKey := createPropertyKey(name)
	if property, ok := isolate.Properties.Lookup(propertyKey); ok {
		if factory, ok := property.(func() Backend); ok {
			return factory(), true
		}
	}
	return nil, false
}

func createPropertyKey(name string) string {
	return "backend." + name
}

func Register(name string, factory func() Backend) {
	isolate.RegisterConfigurator(func(isolate *isolate.Isolate) {
		propertyKey := createPropertyKey(name)
		isolate.Properties.Insert(propertyKey, factory)
	})
}
