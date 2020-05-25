package buildtool

import (
	"github.com/strict-lang/sdk/pkg/buildtool/namespace"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
)

type packageCompilationConfig struct {
	backend backend.Backend
	namespaces *namespace.Table
	outputPath string
}

func compilePackage(config packageCompilationConfig) packageCompilationResult {
	compilation := newPackageCompilation(config)
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
	outputPath string
}

type packageCompilationResult struct {
	diagnostics *diagnostic.Diagnostics
	lineMaps *linemap.Table
}

func newPackageCompilation(config packageCompilationConfig) *packageCompilation {
	return &packageCompilation{
		lineMaps: linemap.NewEmptyTable(),
		backend:  config.backend,
		namespaces:  config.namespaces,
		diagnostics: diagnostic.Empty(),
		outputPath: config.outputPath,
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
	config := compilation.createNamespaceCompilationConfig(namespace)
	diagnostics := compileNamespace(config)
	compilation.addDiagnostics(diagnostics)
	namespace.MarkAsCompiled()
}

func (compilation *packageCompilation) createNamespaceCompilationConfig(
	namespace namespace.Namespace) namespaceCompilationConfig {

	return namespaceCompilationConfig{
		lineMaps:   compilation.lineMaps,
		backend:    compilation.backend,
		namespace:  namespace,
		namespaces: compilation.namespaces,
		outputPath: compilation.outputPath,
	}
}

func (compilation *packageCompilation) addDiagnostics(
	diagnostics *diagnostic.Diagnostics) {

	compilation.diagnostics = compilation.diagnostics.Merge(diagnostics)
}
