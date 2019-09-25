package arduino

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Generation struct {
	backend.Extension

	className string
	parent    *backend.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *syntaxtree.Visitor) {
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

func (generation *Generation) VisitImportStatement(statement *syntaxtree.ImportStatement) {
	generation.parent.EmitFormatted("#include %s", statement.Target.FilePath())
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) VisitClassDeclaration(declaration *syntaxtree.ClassDeclaration) {
	methods, others := extractMethods(declaration.Children)
	generation.writeJustDeclarations(methods)
	fields, setupBody := extractFieldDeclarations(others)
	generation.writeGlobalFieldDeclarations(fields)
	generation.writeMethodDefinitions(methods)
	generation.generateSetupMethod(setupBody)
}

func (generation *Generation) generateSetupMethod(statements []syntaxtree.Node) {
	generation.parent.EmitNode(&syntaxtree.MethodDeclaration{
		Name: &syntaxtree.Identifier{
			Value:        "setup",
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Type: &syntaxtree.ConcreteTypeName{
			Name:         "void",
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Parameters: []*syntaxtree.Parameter{},
		Body: &syntaxtree.BlockStatement{
			Children:     statements,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		NodePosition: syntaxtree.ZeroPosition{},
	})
}

func (generation *Generation) writeGlobalFieldDeclarations(fields []*syntaxtree.FieldDeclaration) {
	for _, field := range fields {
		generation.parent.GenerateFieldDeclaration(field)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeJustDeclarations(methods []*syntaxtree.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitMethodDeclaration(method)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeMethodDefinitions(methods []*syntaxtree.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitNode(method)
		generation.parent.EmitEndOfLine()
	}
}

func extractFieldDeclarations(nodes []syntaxtree.Node) (fields []*syntaxtree.FieldDeclaration, others []syntaxtree.Node) {
	for _, node := range nodes {
		if field, isField := node.(*syntaxtree.FieldDeclaration); isField {
			fields = append(fields, field)
		} else {
			others = append(others, node)
		}
	}
	return
}

func extractMethods(nodes []syntaxtree.Node) (methods []*syntaxtree.MethodDeclaration, others []syntaxtree.Node) {
	for _, node := range nodes {
		if method, isMethod := node.(*syntaxtree.MethodDeclaration); isMethod {
			methods = append(methods, method)
		} else {
			others = append(others, node)
		}
	}
	return
}
