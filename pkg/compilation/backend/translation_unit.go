package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (generation *Generation) GenerateTranslationUnit(unit *syntaxtree2.TranslationUnit) {
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
	generation.Emit("#include <string>\n#include <vector>\n")
}

func (generation *Generation) GenerateMainMethod(nodes []syntaxtree2.Node) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &syntaxtree2.BlockStatement{
		Children: nodes,
	}
	generation.EmitNode(block)
}
