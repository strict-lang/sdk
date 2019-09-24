package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (generation *Generation) GenerateTranslationUnit(unit *syntaxtree.TranslationUnit) {
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

func (generation *Generation) GenerateMainMethod(nodes []syntaxtree.Node) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &syntaxtree.BlockStatement{
		Children: nodes,
	}
	generation.EmitNode(block)
}
