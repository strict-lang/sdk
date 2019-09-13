package backend

import "gitlab.com/strict-lang/sdk/compilation/syntaxtree"

func (generation *Generation) GenerateImportStatement(statement *syntaxtree.ImportStatement) {
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
