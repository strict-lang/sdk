package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"github.com/BenjaminNitschke/Strict/pkg/scope"
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"testing"
)

func TestCodeGeneration(test *testing.T) {
	entry := ast.Method{
		Name: ast.NewIdentifier("commentOnAge"),
		Type: ast.GenericTypeName{
			Name: "list",
			Generic: ast.ConcreteTypeName{
				Name: "number",
			},
		},
		Parameters: []ast.Parameter{
			{
				Name: ast.NewIdentifier("age"),
				Type: ast.ConcreteTypeName{
					Name: "number",
				},
			},
		},
		Body: &ast.BlockStatement{
			Children: []ast.Node{
				&ast.ConditionalStatement{
					Condition: &ast.BinaryExpression{
						LeftOperand:  &ast.Identifier{Value: "age"},
						RightOperand: &ast.NumberLiteral{Value: "18"},
						Operator:     token.SmallerOperator,
					},
					Body: &ast.BlockStatement{
						Children: []ast.Node{
							&ast.ExpressionStatement{
								Expression: &ast.MethodCall{
									Name: ast.Identifier{Value: "logf"},
									Arguments: []ast.Node{
										&ast.StringLiteral{Value: "%d are still young"},
										&ast.Identifier{Value: "age"},
									},
								},
							},
						},
					},
					Else: &ast.BlockStatement{
						Children: []ast.Node{
							&ast.ExpressionStatement{
								Expression: &ast.MethodCall{
									Name: ast.Identifier{Value: "log"},
									Arguments: []ast.Node{
										&ast.StringLiteral{Value: "You are an adult"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	unit := ast.NewTranslationUnit("test", scope.NewRoot(), []ast.Node{&entry,})
	generator := NewCodeGenerator(&unit)
	generator.GenerateMethod(&entry)
	generator.PrintOutput()
}

func (generator *CodeGenerator) PrintOutput() {
	fmt.Println(generator.output.String())
}
