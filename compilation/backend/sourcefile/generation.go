package sourcefile

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
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

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *ast.Visitor) {
	generation.generation = parent
	visitor.VisitClassDeclaration = generation.generateClassDeclaration
}

func (generation *Generation) generateClassDeclaration(declaration *ast.ClassDeclaration) {
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

func splitMethods(nodes []ast.Node) (methods []*ast.MethodDeclaration, remainder []ast.Node) {
	for _, node := range nodes {
		if method, isMethod := node.(*ast.MethodDeclaration); isMethod {
			methods = append(methods, method)
		} else {
			remainder = append(remainder, node)
		}
	}
	return
}

func (generation *Generation) writeMethodDeclaration(declaration *ast.MethodDeclaration) {
	name := fmt.Sprintf("%s::%s", generation.className, declaration.Name.Value)
	instanceMethod := &ast.MethodDeclaration{
		Name: &ast.Identifier{
			Value:        name,
			NodePosition: ast.ZeroPosition{},
		},
		Type:         declaration.Type,
		Parameters:   declaration.Parameters,
		Body:         declaration.Body,
		NodePosition: declaration.NodePosition,
	}
	generation.generation.GenerateMethod(instanceMethod)
}

func (generation *Generation) writeInitMethod(body []ast.Node) {
	generation.writeMethodDeclaration(&ast.MethodDeclaration{
		Name: &ast.Identifier{
			Value:        backend.InitMethodName,
			NodePosition: ast.ZeroPosition{},
		},
		Type: &ast.ConcreteTypeName{
			Name:         "void",
			NodePosition: ast.ZeroPosition{},
		},
		Body: &ast.BlockStatement{
			Children:     body,
			NodePosition: ast.ZeroPosition{},
		},
		Parameters:   []*ast.Parameter{},
		NodePosition: ast.ZeroPosition{},
	})
}
