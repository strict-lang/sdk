package sourcefile

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type Generation struct {
	backend.Extension

	className      string
	generation     *backend.Generation
	hasWrittenInit bool
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *syntaxtree.Visitor) {
	generation.generation = parent
	visitor.VisitClassDeclaration = generation.generateClassDeclaration
}

func (generation *Generation) generateClassDeclaration(declaration *syntaxtree.ClassDeclaration) {
	generation.className = declaration.Name
	methods, initBody := splitMethods(declaration.Children)
	if len(initBody) > 0 {
		generation.writeInitMethod(initBody)
		generation.generation.EmitEndOfLine()
		generation.hasWrittenInit = true
	}
	for _, method := range methods {
		generation.writeMethodDeclaration(method)
		generation.generation.EmitEndOfLine()
	}
}

func createInitBody(members []syntaxtree.Node) (body []syntaxtree.Node) {
	for _, field := range members {
		field, isField := field.(*syntaxtree.FieldDeclaration)
		if !isField {
			continue
		}
		body = append(body, field)
	}
	return
}

func createInitStatement(field *syntaxtree.FieldDeclaration) syntaxtree.Node {
	return &syntaxtree.AssignStatement{
		Target:       field.Name,
		Value:        &syntaxtree.CallExpression{
			Method:       field.TypeName,
			Arguments:    []syntaxtree.Node{},
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Operator:     0,
		NodePosition: nil,
	}
}

func splitMethods(nodes []syntaxtree.Node) (methods []*syntaxtree.MethodDeclaration, remainder []syntaxtree.Node) {
	for _, node := range nodes {
		if method, isMethod := node.(*syntaxtree.MethodDeclaration); isMethod {
			methods = append(methods, method)
		} else {
			remainder = append(remainder, node)
		}
	}
	return
}

func (generation *Generation) writeMethodDeclaration(declaration *syntaxtree.MethodDeclaration) {
	name := fmt.Sprintf("%s::%s", generation.className, declaration.Name.Value)
	instanceMethod := &syntaxtree.MethodDeclaration{
		Name: &syntaxtree.Identifier{
			Value:        name,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Type:         declaration.Type,
		Parameters:   declaration.Parameters,
		Body:         declaration.Body,
		NodePosition: declaration.NodePosition,
	}
	generation.generation.GenerateMethod(instanceMethod)
}

func (generation *Generation) writeInitMethod(body []syntaxtree.Node) {
	generation.writeMethodDeclaration(&syntaxtree.MethodDeclaration{
		Name: &syntaxtree.Identifier{
			Value:        backend.InitMethodName,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Type: &syntaxtree.ConcreteTypeName{
			Name:         "void",
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Body: &syntaxtree.BlockStatement{
			Children:     body,
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Parameters:   []*syntaxtree.Parameter{},
		NodePosition: syntaxtree.ZeroPosition{},
	})
}
