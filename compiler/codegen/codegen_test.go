package codegen

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/scope"
	"gitlab.com/strict-lang/sdk/compiler/token"
	"testing"
)

func TestCodeGeneration(test *testing.T) {
	method := ast.Method{
		Name: ast.NewIdentifier("commentOnAge"),
		Type: ast.ConcreteTypeName{
			Name: "void",
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
									Method: &ast.Identifier{Value: "logf"},
									Arguments: []ast.Node{
										&ast.StringLiteral{Value: "%d is still young"},
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

	unit := ast.NewTranslationUnit("test", scope.NewRoot(), []ast.Node{&method, &call})
	generator := NewCodeGenerator(unit)
	test.Log(generator.Generate())
	// TODO(merlinosayimwen): Validate output
}

func (generator *CodeGenerator) PrintOutput() {
	fmt.Println(generator.output.String())
}
