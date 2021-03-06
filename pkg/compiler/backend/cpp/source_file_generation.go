package cpp

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type SourceFileGeneration struct {
	Extension

	className             string
	generation            *Generation
	hasWrittenInit        bool
	hasWrittenConstructor bool
}

func NewSourceFileGeneration() *SourceFileGeneration {
	return &SourceFileGeneration{}
}

func (generation *SourceFileGeneration) ModifyVisitor(
	parent *Generation, visitor *tree.DelegatingVisitor) {

	generation.generation = parent
	visitor.ClassDeclarationVisitor = generation.generateClassDeclaration
	visitor.MethodDeclarationVisitor = generation.generateMethodDeclaration
	visitor.ConstructorDeclarationVisitor = generation.generateConstructorDeclaration
	generation.importMatchingHeader()
}

func (generation *SourceFileGeneration) importMatchingHeader() {
	className := generation.generation.Unit.Class.Name
	generation.generation.EmitFormatted("#include \"%s.h\"", className)
	generation.generation.EmitEndOfLine()
	generation.generation.EmitEndOfLine()
}

func (generation *SourceFileGeneration) generateClassDeclaration(
	declaration *tree.ClassDeclaration) {

	generation.className = declaration.Name
	members, initBody := filterDeclarations(declaration.Children)
	if len(initBody) > 0 {
		generation.writeInitMethod(ExtractStatements(initBody))
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

func (generation *SourceFileGeneration) writeImplicitConstructor() {
	generation.generateConstructorDeclaration(&tree.ConstructorDeclaration{
		Parameters: []*tree.Parameter{},
		Body:       &tree.StatementBlock{},
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
			Target: &tree.Identifier{
				Value: field.TypeName.BaseName(),
			},
			Arguments: []*tree.CallArgument{},
		},
		Operator: 0,
	}
}

func filterDeclarations(
	nodes []tree.Node) (declarations []tree.Node, remainder []tree.Node) {

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

func (generation *SourceFileGeneration) generateMethodDeclaration(
	declaration *tree.MethodDeclaration) {

	name := fmt.Sprintf("%s::%s", generation.className, declaration.Name.Value)
	instanceMethod := &tree.MethodDeclaration{
		Name: &tree.Identifier{
			Value: name,
		},
		Type:       declaration.Type,
		Parameters: declaration.Parameters,
		Body:       declaration.Body,
	}
	generation.generation.GenerateMethod(instanceMethod)
}

func (generation *SourceFileGeneration) generateConstructorDeclaration(
	declaration *tree.ConstructorDeclaration) {
	generation.hasWrittenConstructor = true
	output := generation.generation
	className := generation.generation.Unit.Class.Name
	output.EmitFormatted("%s::%s", className, className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(" ")
	body := &tree.StatementBlock{
		Children: []tree.Statement{
			generation.generateInitCall(),
			declaration.Body,
		},
	}
	output.EmitNode(body)
}

func (generation *SourceFileGeneration) generateInitCall() tree.Node {
	return &tree.ExpressionStatement{
		Expression: &tree.CallExpression{
			Target: &tree.Identifier{
				Value: InitMethodName,
			},
			Arguments: []*tree.CallArgument{},
		},
	}
}

func (generation *SourceFileGeneration) writeInitMethod(body []tree.Statement) {
	generation.generateMethodDeclaration(&tree.MethodDeclaration{
		Name: &tree.Identifier{
			Value: InitMethodName,
		},
		Type: &tree.ConcreteTypeName{
			Name: "void",
		},
		Body: &tree.StatementBlock{
			Children: body,
		},
		Parameters: []*tree.Parameter{},
	})
}
