package sourcefile

import (
	"fmt"
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Generation struct {
	backend2.Extension

	className      string
	generation     *backend2.Generation
	hasWrittenInit bool
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend2.Generation, visitor *syntaxtree2.Visitor) {
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

func (generation *Generation) generateClassDeclaration(declaration *syntaxtree2.ClassDeclaration) {
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

func createInitBody(members []syntaxtree2.Node) (body []syntaxtree2.Node) {
	for _, field := range members {
		field, isField := field.(*syntaxtree2.FieldDeclaration)
		if !isField {
			continue
		}
		body = append(body, field)
	}
	return
}

func createInitStatement(field *syntaxtree2.FieldDeclaration) syntaxtree2.Node {
	return &syntaxtree2.AssignStatement{
		Target: field.Name,
		Value: &syntaxtree2.CallExpression{
			Method:       field.TypeName,
			Arguments:    []*syntaxtree2.CallArgument{},
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Operator:     0,
		NodePosition: nil,
	}
}

func filterDeclarations(nodes []syntaxtree2.Node) (declarations []syntaxtree2.Node, remainder []syntaxtree2.Node) {
	for _, node := range nodes {
		switch node.(type) {
		case *syntaxtree2.MethodDeclaration, *syntaxtree2.ConstructorDeclaration:
			declarations = append(declarations, node)
			continue
		case *syntaxtree2.FieldDeclaration: // Field declarations are not written
			continue
		default:
			remainder = append(remainder, node)
		}
	}
	return
}

func (generation *Generation) generateMethodDeclaration(declaration *syntaxtree2.MethodDeclaration) {
	name := fmt.Sprintf("%s::%s", generation.className, declaration.Name.Value)
	instanceMethod := &syntaxtree2.MethodDeclaration{
		Name: &syntaxtree2.Identifier{
			Value:        name,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Type:         declaration.Type,
		Parameters:   declaration.Parameters,
		Body:         declaration.Body,
		NodePosition: declaration.NodePosition,
	}
	generation.generation.GenerateMethod(instanceMethod)
}

func (generation *Generation) generateConstructorDeclaration(declaration *syntaxtree2.ConstructorDeclaration) {
	output := generation.generation
	className := generation.generation.Unit.Class.Name
	output.EmitFormatted("%s::%s", className, className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(" ")
	output.EmitNode(declaration.Body)
}

func (generation *Generation) writeInitMethod(body []syntaxtree2.Node) {
	generation.generateMethodDeclaration(&syntaxtree2.MethodDeclaration{
		Name: &syntaxtree2.Identifier{
			Value:        backend2.InitMethodName,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Type: &syntaxtree2.ConcreteTypeName{
			Name:         "void",
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Body: &syntaxtree2.BlockStatement{
			Children:     body,
			NodePosition: syntaxtree2.ZeroPosition{},
		},
		Parameters:   []*syntaxtree2.Parameter{},
		NodePosition: syntaxtree2.ZeroPosition{},
	})
}
