package buildtool

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
)

func compilePackage(
	backend backend.Backend,
	namespaces *namespace.Table) packageCompilationResult {

	compilation := newPackageCompilation(backend, namespaces)
	compilation.run()
	return packageCompilationResult{
		diagnostics: compilation.diagnostics,
		lineMaps:    compilation.lineMaps,
	}
}

type packageCompilation struct {
	lineMaps *linemap.Table
	backend backend.Backend
	namespaces *namespace.Table
	diagnostics *diagnostic.Diagnostics
}

type packageCompilationResult struct {
	diagnostics *diagnostic.Diagnostics
	lineMaps *linemap.Table
}

func newPackageCompilation(
	backend backend.Backend,
	namespaces *namespace.Table) *packageCompilation {

	return &packageCompilation{
		lineMaps: linemap.NewEmptyTable(),
		backend: backend,
		namespaces:  namespaces,
		diagnostics: diagnostic.Empty(),
	}
}

func (compilation *packageCompilation) run() {
	for _, namespace := range compilation.namespaces.List() {
		compilation.compileNamespace(namespace)
	}
}

func (compilation *packageCompilation) compileNamespace(namespace namespace.Namespace) {
	if namespace.IsCompiled() {
		return
	}
	diagnostics := compileNamespace(
		compilation.lineMaps,
		compilation.backend,
		namespace,
		compilation.namespaces)
	compilation.addDiagnostics(diagnostics)
	namespace.MarkAsCompiled()
}

func (compilation *packageCompilation) addDiagnostics(
	diagnostics *diagnostic.Diagnostics) {

	compilation.diagnostics = compilation.diagnostics.Merge(diagnostics)
}
