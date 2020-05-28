package entering

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"strings"
)

func init() {
	passes.Register(newImplicitFactoryPass())
}

const ImplicitFactoryPassId = "GenericResolutionPass"

type ImplicitFactoryPass struct {
	class *tree.ClassDeclaration
	fields []*tree.FieldDeclaration
}

func newImplicitFactoryPass() *ImplicitFactoryPass {
	return &ImplicitFactoryPass{}
}

func (pass *ImplicitFactoryPass) Id() passes.Id {
	return GenericResolutionPassId
}

func (pass *ImplicitFactoryPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ParentAssignPassId)
}

func (pass *ImplicitFactoryPass) Run(context *passes.Context) {
	pass.class = context.Unit.Class
	pass.fields = listFields(context.Unit.Class.Children)
	pass.insertImplicitFactory()
}

func (pass *ImplicitFactoryPass) insertImplicitFactory() {
	pass.class.Children = append(pass.class.Children, pass.createImplicitFactory())
}

func (pass *ImplicitFactoryPass) createImplicitFactory() *tree.MethodDeclaration {
	return &tree.MethodDeclaration{
		Name:       &tree.Identifier{Value: "factory"},
		Parameters: pass.createFactoryParameters(),
		Body:       &tree.StatementBlock{
			Children: pass.createAssignStatements(),
		},
		Factory:    true,
	}
}

func (pass *ImplicitFactoryPass) createFactoryParameters() (parameters []*tree.Parameter) {
	for _, field := range pass.fields {
		parameter := &tree.Parameter{
			Type:   field.TypeName,
			Name:   &tree.Identifier{Value: strings.ToLower(field.Name.Value)},
		}
		parameters = append(parameters, parameter)
	}
	return
}

func (pass *ImplicitFactoryPass) createAssignStatements() (statements []tree.Statement) {
	for _, field := range pass.fields {
		statement := &tree.AssignStatement{
			Target:    &tree.ChainExpression{
				Expressions: []tree.Expression{
					&tree.Identifier{Value: `this`},
					&tree.Identifier{Value: field.Name.Value},
				},
			},
			Value:    &tree.Identifier{Value: strings.ToLower(field.Name.Value)},
			Operator: token.AssignOperator,
		}
		statements = append(statements, statement)
	}
	return
}

func listFields(elements []tree.Node) (fields []*tree.FieldDeclaration) {
	for _, element := range elements {
		if field, isField := element.(*tree.FieldDeclaration); isField {
			fields = append(fields, field)
		}
	}
	return
}
