package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

func (generation *Generation) GenerateTranslationUnit(unit *ast.TranslationUnit) {
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
	generation.Emit("#include <string>\n")
}

func (generation *Generation) GenerateMainMethod(nodes []ast.Node) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &ast.BlockStatement{
		Children: nodes,
	}
	generation.EmitNode(block)
}
