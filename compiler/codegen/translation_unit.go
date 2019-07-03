package codegen

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
)

func (generator *CodeGenerator) GenerateTranslationUnit(unit *ast.TranslationUnit) {
	generator.Emit("// Code is generated by the strict compiler.\n")
	generator.Emit("\n#include <strict/all.hh>\n\n")
	methods, others := splitTopLevelNodes(unit)
	for _, method := range methods {
		method.Accept(generator.generators)
	}
	generator.GenerateMainMethod(others)

	generator.Emit("\n")
}

func splitTopLevelNodes(unit *ast.TranslationUnit) (methods []ast.Node, others []ast.Node) {
	for _, node := range unit.Children {
		if _, ok := node.(*ast.Method); ok {
			methods = append(methods, node)
		} else {
			others = append(others, node)
		}
	}
	return methods, others
}

func (generator *CodeGenerator) GenerateMainMethod(nodes []ast.Node) {
	generator.Emit("int main(int argc, char **argv) ")
	block := &ast.BlockStatement{
		Children: nodes,
	}
	block.Accept(generator.generators)
}
