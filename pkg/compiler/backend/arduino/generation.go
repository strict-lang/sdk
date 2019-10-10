package arduino

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Generation struct {
	backend.Extension

	className string
	parent    *backend.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *tree.Visitor) {
	generation.parent = parent
	parent.DisableNamespaceSelectors()
	parent.DisableStdlibClassImport()
	visitor.VisitClassDeclaration = generation.VisitClassDeclaration
	visitor.VisitImportStatement = generation.VisitImportStatement
	generation.importArduinoHeader()
}

func (generation *Generation) importArduinoHeader() {
	generation.parent.Emit(`#include "arduino.h"`)
	generation.parent.EmitEndOfLine()
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) VisitImportStatement(statement *tree.ImportStatement) {
	generation.parent.EmitFormatted("#include %s", statement.Target.FilePath())
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) VisitClassDeclaration(declaration *tree.ClassDeclaration) {
	methods, others := extractMethods(declaration.Children)
	generation.writeJustDeclarations(methods)
	fields, setupBody := extractFieldDeclarations(others)
	generation.writeGlobalFieldDeclarations(fields)
	generation.writeMethodDefinitions(methods)
	generation.generateSetupMethod(setupBody)
}

func (generation *Generation) generateSetupMethod(statements []tree.Node) {
	generation.parent.EmitNode(&tree.MethodDeclaration{
		Name: &tree.Identifier{
			Value:        "setup",
			NodePosition: tree.ZeroArea{},
		},
		Type: &tree.ConcreteTypeName{
			Name:         "void",
			NodePosition: tree.ZeroArea{},
		},
		Parameters: []*tree.Parameter{},
		Body: &tree.BlockStatement{
			Children:     statements,
			NodePosition: tree.ZeroArea{},
		},
		NodePosition: tree.ZeroArea{},
	})
}

func (generation *Generation) writeGlobalFieldDeclarations(fields []*tree.FieldDeclaration) {
	for _, field := range fields {
		generation.parent.GenerateFieldDeclaration(field)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeJustDeclarations(methods []*tree.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitMethodDeclaration(method)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeMethodDefinitions(methods []*tree.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitNode(method)
		generation.parent.EmitEndOfLine()
	}
}

func extractFieldDeclarations(nodes []tree.Node) (fields []*tree.FieldDeclaration, others []tree.Node) {
	for _, node := range nodes {
		if field, isField := node.(*tree.FieldDeclaration); isField {
			fields = append(fields, field)
		} else {
			others = append(others, node)
		}
	}
	return
}

func extractMethods(nodes []tree.Node) (methods []*tree.MethodDeclaration, others []tree.Node) {
	for _, node := range nodes {
		if method, isMethod := node.(*tree.MethodDeclaration); isMethod {
			methods = append(methods, method)
		} else {
			others = append(others, node)
		}
	}
	return
}
