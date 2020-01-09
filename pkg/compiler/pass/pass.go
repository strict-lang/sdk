package pass

import "gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"

type Context struct {
	Unit *tree.TranslationUnit
}

type Id string

type Pass struct {
	Run func(context *Context)
	Dependencies []*Pass
}

type Builder struct {
	dependencies []*Pass
	runner func(context *Context)
}

func NewPass() *Builder {
	return &Builder{}
}

func (builder *Builder) AddDependency(dependency *Pass) *Builder {
	builder.dependencies = append(builder.dependencies, dependency)
	return builder
}

func (builder *Builder) RunWith(runner func(context *Context)) *Builder {
	builder.runner = runner
	return builder
}

func (builder *Builder) Create() *Pass {
	copiedDependencies := make([]*Pass, len(builder.dependencies))
	copy(copiedDependencies, builder.dependencies)
	return &Pass{
		Run: builder.runner,
		Dependencies: copiedDependencies,
	}
}