package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"path/filepath"
	"strings"
)

func (generation *Generation) GenerateImportStatement(statement *tree.ImportStatement) {
	moduleName := statement.ModuleName()
	path := generation.resolveImportPath(statement.Target.Namespace())
	generation.importModule(moduleName, path)
}

const namespaceHeader = "_namespace.h"

func (generation *Generation) resolveImportPath(namespace string) string {
	basePath := strings.ReplaceAll(namespace, ".", string(filepath.Separator))
	return filepath.Join(basePath, namespaceHeader)
}

func (generation *Generation) importModule(name, path string) {
	generation.importModules[name] = path
	generation.includeIntoNamespace(name, path)
}

func (generation *Generation) includeIntoNamespace(name, path string) {
	generation.EmitFormatted("namespace %s {\n  #include %s\n}\n", name, path)
}
