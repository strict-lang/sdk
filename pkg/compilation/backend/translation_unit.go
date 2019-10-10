package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/syntax/tree"
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

func (generation *Generation) GenerateMainMethod(nodes []tree.Node) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &tree.BlockStatement{
		Children: nodes,
	}
	generation.EmitNode(block)
}
