package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

func (generator *CodeGenerator) GenerateTranslationUnit(unit *ast.TranslationUnit) {
	methods, nonMethods := splitTopLevelNodes(unit)
	sharedVariableDeclarations, others := splitSharedVariableDeclarations(nonMethods)
	for _, declaration := range sharedVariableDeclarations {
		generator.EmitNode(declaration)
	}
	for _, method := range methods {
		generator.EmitNode(method)
	}
	generator.GenerateMainMethod(others)
	generator.Emit("\n")
}

func splitSharedVariableDeclarations(nodes []ast.Node) (declarations []ast.Node, others []ast.Node) {
	for _, node := range nodes {
		if _, ok := node.(*ast.SharedVariableDeclaration); ok {
			declarations = append(declarations, node)
		} else {
			others = append(others, node)
		}
	}
	return declarations, others
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
	generator.EmitNode(block)
}
