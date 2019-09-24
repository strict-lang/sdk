package arduino

import (
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Generation struct {
	backend2.Extension

	className string
	parent    *backend2.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend2.Generation, visitor *syntaxtree2.Visitor) {
	generation.parent = parent
	parent.DisableNamespaceSelectors()
	visitor.VisitClassDeclaration = generation.VisitClassDeclaration
	visitor.VisitImportStatement = generation.VisitImportStatement
	generation.importArduinoHeader()
}

func (generation *Generation) importArduinoHeader() {
	generation.parent.Emit(`#include "arduino.h"`)
	generation.parent.EmitEndOfLine()
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) VisitImportStatement(statement *syntaxtree2.ImportStatement) {
	generation.parent.EmitFormatted("#include %s", statement.Target.FilePath())
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) VisitClassDeclaration(declaration *syntaxtree2.ClassDeclaration) {
	methods, others := extractMethods(declaration.Children)
	generation.writeJustDeclarations(methods)
	fields, setupBody := extractFieldDeclarations(others)
	generation.writeGlobalFieldDeclarations(fields)
	generation.writeMethodDefinitions(methods)
	generation.generateSetupMethod(setupBody)
}

func (generation *Generation) generateSetupMethod(statements []syntaxtree2.Node) {
	generation.parent.EmitNode(&syntaxtree2.MethodDeclaration{
		Name: &syntaxtree2.Identifier{
			Value:        "setup",
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Type: &syntaxtree2.ConcreteTypeName{
			Name:         "void",
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Parameters: []*syntaxtree2.Parameter{},
		Body: &syntaxtree2.BlockStatement{
			Children:     statements,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		NodePosition: syntaxtree2.ZeroPosition{},
	})
}

func (generation *Generation) writeGlobalFieldDeclarations(fields []*syntaxtree2.FieldDeclaration) {
	for _, field := range fields {
		generation.parent.GenerateFieldDeclaration(field)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeJustDeclarations(methods []*syntaxtree2.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitMethodDeclaration(method)
		generation.parent.Emit(";")
		generation.parent.EmitEndOfLine()
	}
	generation.parent.EmitEndOfLine()
}

func (generation *Generation) writeMethodDefinitions(methods []*syntaxtree2.MethodDeclaration) {
	for _, method := range methods {
		generation.parent.EmitNode(method)
		generation.parent.EmitEndOfLine()
	}
}

func extractFieldDeclarations(nodes []syntaxtree2.Node) (fields []*syntaxtree2.FieldDeclaration, others []syntaxtree2.Node) {
	for _, node := range nodes {
		if field, isField := node.(*syntaxtree2.FieldDeclaration); isField {
			fields = append(fields, field)
		} else {
			others = append(others, node)
		}
	}
	return
}

func extractMethods(nodes []syntaxtree2.Node) (methods []*syntaxtree2.MethodDeclaration, others []syntaxtree2.Node) {
	for _, node := range nodes {
		if method, isMethod := node.(*syntaxtree2.MethodDeclaration); isMethod {
			methods = append(methods, method)
		} else {
			others = append(others, node)
		}
	}
	return
}
