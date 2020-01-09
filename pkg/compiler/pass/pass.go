package pass

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Context struct {
	Unit *tree.TranslationUnit
	Diagnostic *diagnostic.Bag
}

type Id string

type Pass interface {
	Run(context *Context)
	Dependencies() Set
}

type Set []Pass

var EmptySet = Set{}