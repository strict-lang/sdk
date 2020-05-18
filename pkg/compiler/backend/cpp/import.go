package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (generation *Generation) GenerateImportStatement(statement *tree.ImportStatement) {
	moduleName := statement.ModuleName()
	generation.importModule(moduleName, statement.Target.Namespace())
}

func (generation *Generation) importModule(name, path string) {
	generation.importModules[name] = path
	generation.includeIntoNamespace(name, path)
}

func (generation *Generation) includeIntoNamespace(name, path string) {
	generation.EmitFormatted("namespace %s {\n  #include %s\n}\n", name, path)
}
