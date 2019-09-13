package sourcefile

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/backend"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
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
	visitor.VisitMethodDeclaration = generation.generateMethodDeclaration
	visitor.VisitConstructorDeclaration = generation.generateConstructorDeclaration
	generation.importMatchingHeader()
}

func (generation *Generation) importMatchingHeader() {
	className := generation.generation.Unit.Class.Name
	generation.generation.EmitFormatted("#include \"%s.h\"", className)
	generation.generation.EmitEndOfLine()
	generation.generation.EmitEndOfLine()
}

func (generation *Generation) generateClassDeclaration(declaration *syntaxtree.ClassDeclaration) {
	generation.className = declaration.Name
	members, initBody := filterDeclarations(declaration.Children)
	if len(initBody) > 0 {
		generation.writeInitMethod(initBody)
		generation.generation.EmitEndOfLine()
		generation.hasWrittenInit = true
	}
	for _, member := range members {
		generation.generation.EmitNode(member)
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
		Target: field.Name,
		Value: &syntaxtree.CallExpression{
			Method:       field.TypeName,
			Arguments:    []*syntaxtree.CallArgument{},
			NodePosition: syntaxtree.ZeroPosition{},
		},
		Operator:     0,
		NodePosition: nil,
	}
}

func filterDeclarations(nodes []syntaxtree.Node) (declarations []syntaxtree.Node, remainder []syntaxtree.Node) {
	for _, node := range nodes {
		switch node.(type) {
		case *syntaxtree.MethodDeclaration, *syntaxtree.ConstructorDeclaration:
			declarations = append(declarations, node)
			continue
		case *syntaxtree.FieldDeclaration: // Field declarations are not written
			continue
		default:
			remainder = append(remainder, node)
		}
	}
	return
}

func (generation *Generation) generateMethodDeclaration(declaration *syntaxtree.MethodDeclaration) {
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

func (generation *Generation) generateConstructorDeclaration(declaration *syntaxtree.ConstructorDeclaration) {
	output := generation.generation
	className := generation.generation.Unit.Class.Name
	output.EmitFormatted("%s::%s", className, className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(" ")
	output.EmitNode(declaration.Body)
}

func (generation *Generation) writeInitMethod(body []syntaxtree.Node) {
	generation.generateMethodDeclaration(&syntaxtree.MethodDeclaration{
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
