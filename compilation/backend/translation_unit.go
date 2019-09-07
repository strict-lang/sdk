package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

func (generation *Generation) GenerateTranslationUnit(unit *ast.TranslationUnit) {
	importStatements, nonImports := splitImportStatements(unit.Children)
	methods, nonMethods := splitMethodDeclarations(nonImports)
	generation.generateImplicitImports()
	generation.generateSection(importStatements)
	generation.generateSection(methods)
	generation.GenerateMainMethod(nonMethods)
	generation.Emit("\n")
}

func (generation *Generation) generateImplicitImports() {
	generation.Emit("#include <string>\n")
}

func (generation *Generation) generateSection(nodes []ast.Node) {
	generation.generateAll(nodes)
	if len(nodes) > 0 {
		generation.Emit("\n")
	}
}

func (generation *Generation) generateAll(nodes []ast.Node) {
	for _, node := range nodes {
		generation.EmitNode(node)
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

func (generation *Generation) GenerateMainMethod(nodes []ast.Node) {
	generation.Emit("int main(int argc, char **argv) ")
	block := &ast.BlockStatement{
		Children: nodes,
	}
	generation.EmitNode(block)
}
