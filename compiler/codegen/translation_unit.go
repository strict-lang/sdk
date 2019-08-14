package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

func (generator *CodeGenerator) GenerateTranslationUnit(unit *ast.TranslationUnit) {
	methods, nonMethods := splitTopLevelNodes(unit)
	importStatements, nonImports := splitImportStatements(nonMethods)
	sharedVariableDeclarations, others := splitSharedVariableDeclarations(nonImports)
	generator.generateSection(importStatements)
	generator.generateSection(sharedVariableDeclarations)
	generator.generateSection(methods)
	generator.GenerateMainMethod(others)
	generator.Emit("\n")
}

func (generator *CodeGenerator) generateSection(nodes []ast.Node) {
	generator.generateAll(nodes)
	if len(nodes) > 0 {
		generator.Emit("\n")
	}
}

func (generator *CodeGenerator) generateAll(nodes []ast.Node) {
	for _, node := range nodes {
		generator.EmitNode(node)
	}
}

func splitImportStatements(nodes []ast.Node) (importStatements []ast.Node, others []ast.Node) {
	for _, node := range nodes {
		if _, ok := node.(*ast.ImportStatement); ok {
			importStatements = append(importStatements, node)
		} else {
			others = append(others, node)
		}
	}
	return
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
	if generator.settings.IsTargetingArduino {
		generator.Emit("void setup() ")
	} else {
		generator.Emit("int main(int argc, char **argv) ")
	}
	block := &ast.BlockStatement{
		Children: nodes,
	}
	generator.EmitNode(block)
}
