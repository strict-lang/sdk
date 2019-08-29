package backend

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/scope"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"testing"
)

func TestCodeGeneration(test *testing.T) {
	method := ast.MethodDeclaration{
		Name: &ast.Identifier{Value: "commentOnAge"},
		Type: &ast.ConcreteTypeName{
			Name: "void",
		},
		Parameters: []*ast.Parameter{
			{
				Name: &ast.Identifier{
				 	Value: "age",
				},
				Type: &ast.ConcreteTypeName{
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
					Consequence: &ast.BlockStatement{
						Children: []ast.Node{
							&ast.ExpressionStatement{
								Expression: &ast.MethodCall{
									Method: &ast.Identifier{Value: "logf"},
									Arguments: []ast.Node{
										&ast.StringLiteral{Value: "%d is still young"},
										&ast.Identifier{Value: "age"},
									},
								},
							},
						},
					},
					Alternative: &ast.BlockStatement{
						Children: []ast.Node{
							&ast.ExpressionStatement{
								Expression: &ast.MethodCall{
									Method: &ast.Identifier{Value: "log"},
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

	call := ast.ExpressionStatement{
		Expression: &ast.MethodCall{
			Method: &ast.Identifier{Value: "commentOnAge"},
			Arguments: []ast.Node{
				&ast.MethodCall{
					Method: &ast.Identifier{Value: "inputNumber"},
					Arguments: []ast.Node{
						&ast.StringLiteral{
							Value: "How old are you?",
						},
					},
				},
			},
		},
	}

	unit := ast.NewTranslationUnit("test", scope.NewRootScope(), []ast.Node{&method, &call})
	generator := NewGeneration(Settings{
		IsTargetingArduino: false,
	},
	unit)
	test.Log(generator.Generate())
	// TODO(merlinosayimwen): Validate output
}

func (generation *Generation) PrintOutput() {
	fmt.Println(generation.output.String())
}
