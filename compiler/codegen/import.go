package codegen

import "gitlab.com/strict-lang/sdk/compiler/ast"

func (generator *CodeGenerator) GenerateImportStatement(statement *ast.ImportStatement) {
	moduleName := statement.ModuleName()
	generator.importModule(moduleName, statement.Path)
}

func (generator *CodeGenerator) importModule(name, path string) {
	generator.importModules[name] = path
	generator.includeIntoNamespace(name, path)
}

func (generator *CodeGenerator) includeIntoNamespace(name, path string) {
	generator.Emitf("namespace %s {\n  #include <%s>\n}\n", name, path)
}
