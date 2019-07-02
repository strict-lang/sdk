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
		Name: ast.NewIdentifier("divisibleNumbers"),
		Type: ast.GenericTypeName{
			Name: "list",
			Generic: ast.ConcreteTypeName{
				Name:"number",
			},
		},
		Parameters: []ast.Parameter{
			{
				Name: ast.NewIdentifier("limit"),
				Type: ast.ConcreteTypeName{
					Name: "number",
				},
			},
		},
		Body: &ast.BlockStatement{
			Children: []ast.Node{
				&ast.ConditionalStatement{
					Condition: &ast.BinaryExpression{
						LeftOperand: &ast.Identifier{Value: "limit"},
						RightOperand: &ast.NumberLiteral{Value: "10"},
						Operator: token.SmallerEqualsOperator,
					},
					Body: &ast.BlockStatement{
						Children: []ast.Node{
							&ast.MethodCall{
								Name: ast.Identifier{Value: "logf"},
								Arguments: []ast.Node{
									&ast.StringLiteral{Value: "The limit of %d is smaller than 11"},
									&ast.Identifier{Value: "limit"},
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