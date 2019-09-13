package backend

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"testing"
)

func TestCodeGeneration(test *testing.T) {
	method := syntaxtree.MethodDeclaration{
		Name: &syntaxtree.Identifier{Value: "commentOnAge"},
		Type: &syntaxtree.ConcreteTypeName{
			Name: "void",
		},
		Parameters: []*syntaxtree.Parameter{
			{
				Name: &syntaxtree.Identifier{
					Value: "age",
				},
				Type: &syntaxtree.ConcreteTypeName{
					Name: "number",
				},
			},
		},
		Body: &syntaxtree.BlockStatement{
			Children: []syntaxtree.Node{
				&syntaxtree.ConditionalStatement{
					Condition: &syntaxtree.BinaryExpression{
						LeftOperand:  &syntaxtree.Identifier{Value: "age"},
						RightOperand: &syntaxtree.NumberLiteral{Value: "18"},
						Operator:     token.SmallerOperator,
					},
					Consequence: &syntaxtree.BlockStatement{
						Children: []syntaxtree.Node{
							&syntaxtree.ExpressionStatement{
								Expression: &syntaxtree.CallExpression{
									Method: &syntaxtree.Identifier{Value: "logf"},
									Arguments: []syntaxtree.Node{
										&syntaxtree.StringLiteral{Value: "%d is still young"},
										&syntaxtree.Identifier{Value: "age"},
									},
								},
							},
						},
					},
					Alternative: &syntaxtree.BlockStatement{
						Children: []syntaxtree.Node{
							&syntaxtree.ExpressionStatement{
								Expression: &syntaxtree.CallExpression{
									Method: &syntaxtree.Identifier{Value: "log"},
									Arguments: []syntaxtree.Node{
										&syntaxtree.StringLiteral{Value: "You are an adult"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	call := syntaxtree.ExpressionStatement{
		Expression: &syntaxtree.CallExpression{
			Method: &syntaxtree.Identifier{Value: "commentOnAge"},
			Arguments: []syntaxtree.Node{
				&syntaxtree.CallExpression{
					Method: &syntaxtree.Identifier{Value: "inputNumber"},
					Arguments: []syntaxtree.Node{
						&syntaxtree.StringLiteral{
							Value: "How old are you?",
						},
					},
				},
			},
		},
	}

	unit := &syntaxtree.TranslationUnit{
		Name:    "test",
		Imports: []*syntaxtree.ImportStatement{},
		Class: &syntaxtree.ClassDeclaration{
			Name:         "test",
			Parameters:   []syntaxtree.ClassParameter{},
			SuperTypes:   []syntaxtree.TypeName{},
			Children:     []syntaxtree.Node{&method, &call},
			NodePosition: syntaxtree.ZeroPosition{},
		},
		NodePosition: syntaxtree.ZeroPosition{},
	}
	generator := NewGeneration(unit)
	test.Log(generator.Generate())
}

func (generation *Generation) PrintOutput() {
	fmt.Println(generation.output.String())
}
