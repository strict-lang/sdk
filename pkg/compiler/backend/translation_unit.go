package backend

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (generation *Generation) GenerateTranslationUnit(unit *tree.TranslationUnit) {
	generation.generateImplicitImports()
	for _, importStatement := range unit.Imports {
		generation.EmitNode(importStatement)
		generation.EmitEndOfLine()
	}
	generation.EmitEndOfLine()
	generation.EmitNode(unit.Class)
	generation.Emit("\n")
}

func (generation *Generation) generateImplicitImports() {
	if generation.shouldImportStdlibClasses {
		generation.Emit("#include <string>\n#include <vector>\n")
	}
}

func (generation *Generation) GenerateMainMethod(nodes []tree.Statement) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &tree.StatementBlock{
		Children: nodes,
	}
	generation.EmitNode(block)
}
