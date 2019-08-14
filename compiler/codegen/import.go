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
	if statement.Alias == nil || statement.Alias.Value == "" {
		return moduleNameByPath(statement.Path)
	}
	return statement.Alias.Value
}

func moduleNameByPath(path string) string {
	var begin = 0
	if strings.Contains(path, "/") {
		begin = strings.LastIndex(path, "/") + 1
	}
	var end int
	if strings.HasSuffix(path, ".h") {
		end = len(path) - 2
	} else {
		end = len(path)
	}
	return path[begin:end]
}

func (generator *CodeGenerator) importModule(name, path string) {
	generator.importModules[name] = path
	generator.includeIntoNamespace(name, path)
}

func (generator *CodeGenerator) includeIntoNamespace(name, path string) {
	generator.Emitf("namespace %s {\n  #include <%s>\n}\n", name, path)
}