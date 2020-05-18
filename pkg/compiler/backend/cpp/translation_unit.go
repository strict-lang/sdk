package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
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

	if generation.isGeneratingApp {
		generation.generateMainMethodCall()
	}
}

func (generation *Generation) generateMainMethodCall() {
	generation.Emit(`void main(int argc, char** argv) {`)
	generation.IncreaseIndent()
	generation.EmitIndent()
	generation.EmitFormatted(`auto app = new %s();`, generation.Unit.Class.Name)
	generation.EmitIndent()
	generation.Emit("app.Run();")
	generation.EmitEndOfLine()
	generation.DecreaseIndent()
	generation.Emit(`}`)
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
