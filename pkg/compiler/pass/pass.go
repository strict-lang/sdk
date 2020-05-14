package pass

import (
	"errors"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
)

type Context struct {
	Unit       *tree.TranslationUnit
	Diagnostic *diagnostic.Bag
	Isolate    *isolate.Isolate
}

type Id string

type Pass interface {
	Id() Id
	Run(context *Context)
	Dependencies(isolate *isolate.Isolate) Set
}

type Set []Pass

var EmptySet = Set{}

func ListInIsolate(isolate *isolate.Isolate, names ...string) Set {
	var passes []Pass
	for _, name := range names {
		if pass, ok := findPassInProperties(name, isolate.Properties); ok {
			passes = append(passes, pass)
		}
	}
	return passes
}

func findPassInProperties(
	name string, table *isolate.ThreadLocalPropertyTable) (Pass, bool) {

	if value, ok := table.Lookup(name); ok {
		pass, isPass := value.(Pass)
		return pass, isPass
	}
	return nil, false
}

func Register(pass Pass) {
	isolate.RegisterConfigurator(func(isolate *isolate.Isolate) {
		name := string(pass.Id())
		isolate.Properties.Insert(name, pass)
	})
}

func RunWithId(id Id, context *Context) error {
	execution, ok := NewExecution(id, context)
	if !ok {
		return errors.New("could not create passes")
	}
	return execution.Run()
}