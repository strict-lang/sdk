package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/scope"
	"github.com/BenjaminNitschke/Strict/compiler/token"
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
									Name: ast.Identifier{Value: "logf"},
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

	call := ast.ExpressionStatement{
		Expression: &ast.MethodCall{
			Name: ast.Identifier{Value: "commentOnAge"},
			Arguments: []ast.Node{
				&ast.MethodCall{
					Name: ast.Identifier{Value: "inputNumber"},
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
	generator.Generate()
	// TODO(merlinosayimwen): Validate ouput
}

func (generator *CodeGenerator) PrintOutput() {
	fmt.Println(generator.output.String())
}
