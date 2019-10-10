package sourcefile

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	 "gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Generation struct {
	backend.Extension

	className      string
	generation     *backend.Generation
	hasWrittenInit bool
	hasWrittenConstructor bool
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *tree.Visitor) {
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

func (generation *Generation) generateClassDeclaration(declaration *tree.ClassDeclaration) {
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
	if !generation.hasWrittenConstructor {
		generation.writeImplicitConstructor()
	}
}

func (generation *Generation) writeImplicitConstructor() {
	generation.generateConstructorDeclaration(&tree.ConstructorDeclaration{
		Parameters:   []*tree.Parameter{},
		Body:         &tree.BlockStatement{},
		NodePosition: tree.ZeroArea{},
	})
}

func createInitBody(members []tree.Node) (body []tree.Node) {
	for _, field := range members {
		field, isField := field.(*tree.FieldDeclaration)
		if !isField {
			continue
		}
		body = append(body, field)
	}
	return
}

func createInitStatement(field *tree.FieldDeclaration) tree.Node {
	return &tree.AssignStatement{
		Target: field.Name,
		Value: &tree.CallExpression{
			Target:       field.TypeName,
			Arguments:    []*tree.CallArgument{},
			NodePosition: tree.ZeroArea{},
		},
		Operator:     0,
		NodePosition: nil,
	}
}

func filterDeclarations(nodes []tree.Node) (declarations []tree.Node, remainder []tree.Node) {
	for _, node := range nodes {
		switch node.(type) {
		case *tree.MethodDeclaration, *tree.ConstructorDeclaration:
			declarations = append(declarations, node)
			continue
		case *tree.FieldDeclaration: // Field declarations are not written
			continue
		default:
			remainder = append(remainder, node)
		}
	}
	return
}

func (generation *Generation) generateMethodDeclaration(declaration *tree.MethodDeclaration) {
	name := fmt.Sprintf("%s::%s", generation.className, declaration.Name.Value)
	instanceMethod := &tree.MethodDeclaration{
		Name: &tree.Identifier{
			Value:        name,
			NodePosition: tree.ZeroArea{},
		},
		Type:         declaration.Type,
		Parameters:   declaration.Parameters,
		Body:         declaration.Body,
		NodePosition: declaration.NodePosition,
	}
	generation.generation.GenerateMethod(instanceMethod)
}

func (generation *Generation) generateConstructorDeclaration(declaration *tree.ConstructorDeclaration) {
	generation.hasWrittenConstructor = true
	output := generation.generation
	className := generation.generation.Unit.Class.Name
	output.EmitFormatted("%s::%s", className, className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(" ")
	body := &tree.BlockStatement{
		Children:     []tree.Node{
			generation.generateInitCall(),
			declaration.Body,
		}	,
		NodePosition: tree.ZeroArea{},
	}
	output.EmitNode(body)
}

func (generation *Generation) generateInitCall() tree.Node {
	return &tree.ExpressionStatement{
		Expression: &tree.CallExpression{
		Target:       &tree.Identifier{
			Value:        backend.InitMethodName,
			NodePosition: tree.ZeroArea{},
		},
		Arguments:    []*tree.CallArgument{},
		NodePosition: tree.ZeroArea{},
	},
	}
}

func (generation *Generation) writeInitMethod(body []tree.Node) {
	generation.generateMethodDeclaration(&tree.MethodDeclaration{
		Name: &tree.Identifier{
			Value:        backend.InitMethodName,
			NodePosition: tree.ZeroArea{},
		},
		Type: &tree.ConcreteTypeName{
			Name:         "void",
			NodePosition: tree.ZeroArea{},
		},
		Body: &tree.BlockStatement{
			Children:     body,
			NodePosition: tree.ZeroArea{},
		},
		Parameters:   []*tree.Parameter{},
		NodePosition: tree.ZeroArea{},
	})
}
