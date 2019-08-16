package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

func (generator *CodeGenerator) GenerateTranslationUnit(unit *ast.TranslationUnit) {
	importStatements, nonImports := splitImportStatements(unit.Children)
	methods, nonMethods := splitMethodDeclarations(nonImports)
	generator.generateSection(importStatements)
	generator.generateSection(methods)
	generator.GenerateMainMethod(nonMethods)
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

func splitMethodDeclarations(nodes []ast.Node) (declarations []ast.Node, others []ast.Node) {
	for _, node := range nodes {
		if _, ok := node.(*ast.MethodDeclaration); ok {
			declarations = append(declarations, node)
		} else {
			others = append(others, node)
		}
	}
	return
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
