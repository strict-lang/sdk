package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (generation *Generation) GenerateImportStatement(statement *syntaxtree2.ImportStatement) {
	moduleName := statement.ModuleName()
	generation.importModule(moduleName, statement.Target.FilePath())
}

func (generation *Generation) importModule(name, path string) {
	generation.importModules[name] = path
	generation.includeIntoNamespace(name, path)
}

func (generation *Generation) includeIntoNamespace(name, path string) {
	generation.EmitFormatted("namespace %s {\n  #include %s\n}\n", name, path)
}
