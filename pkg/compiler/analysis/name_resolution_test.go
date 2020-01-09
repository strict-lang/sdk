package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"testing"
)

var testCode = &tree.TranslationUnit{
	Name:    "Test",
	Imports: []*tree.ImportStatement{},
	Class:   &tree.ClassDeclaration{
		Name:       "Test",
		Children:   []tree.Node{
			&tree.MethodDeclaration{
				Name:       &tree.Identifier{Value: "add"},
				Type:       &tree.ConcreteTypeName{Name: "int"},
				Parameters: tree.ParameterList{
					&tree.Parameter{
						Type:   &tree.ConcreteTypeName{Name: "int"},
						Name:   &tree.Identifier{Value: "left"},
					},
					&tree.Parameter{
						Type:   &tree.ConcreteTypeName{Name: "int"},
						Name:   &tree.Identifier{Value: "right"},
					},
				},
				Body:       &tree.StatementBlock{
					Children: []tree.Statement{
						&tree.ReturnStatement{
							Value:  &tree.BinaryExpression{
								LeftOperand:  &tree.Identifier{Value:  "left"},
								RightOperand: &tree.Identifier{Value: "right"},
								Operator:     token.AddOperator,
							},
						},
					},
				},
			},
		},
	},
}

func TestNameResolutionPass(testing *testing.T) {
	context := &passes.Context{
		Unit:       testCode,
		Diagnostic: diagnostic.NewBag(),
		Isolate:    isolate.SingleThreaded(),
	}
	execution := passes.NewExecution(&NameResolutionPass{}, context)
	if err := execution.Run(); err != nil {
		testing.Error(err)
	}
}
