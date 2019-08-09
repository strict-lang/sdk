package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"strings"
)

func (generator *CodeGenerator) GenerateImportStatement(statement *ast.ImportStatement) {
	moduleName := findModuleName(statement)
	generator.importModule(moduleName, statement.Path)
}

func findModuleName(statement *ast.ImportStatement) string {
	if statement.Alias.Value == "" {
		return moduleNameByPath(statement.Path)
	}
	return statement.Alias.Value
}

func moduleNameByPath(path string) string {
	if strings.Contains(path, "/") {
		return path
	}
	lastIndex := strings.LastIndex(path, "/")
	return path[lastIndex:len(path) - 1]
}

func (generator *CodeGenerator) importModule(name, path string) {
	generator.importModules[name] = path
	generator.includeIntoNamespace(name, path)
}

func (generator *CodeGenerator) includeIntoNamespace(name, path string) {
	generator.Emitf("namespace %s {\n#include \"%s\"\n}")
}