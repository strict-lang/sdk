package syntax

import (
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strings"
	"testing"
)

func TestParsing_ParseMethodDeclaration(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `
method range(begin int, end int) returns List<int>
  for number from begin to end
    yield number
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name: &tree.Identifier{Value: `range`},
					Type: &tree.GenericTypeName{
						Name:    "List",
						Generic: &tree.ConcreteTypeName{Name: `int`},
					},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `int`},
							Name: &tree.Identifier{Value: `begin`},
						},
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `int`},
							Name: &tree.Identifier{Value: `end`},
						},
					},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.RangedLoopStatement{
								Field: &tree.Identifier{Value: `number`},
								Begin: &tree.Identifier{Value: `begin`},
								End:   &tree.Identifier{Value: `end`},
								Body: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.YieldStatement{
											Value: &tree.Identifier{Value: `number`},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Input: `
method printList(numbers List<int>)
  for number in numbers
    printf("%d ", number)
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name: &tree.Identifier{Value: `printList`},
					Type: &tree.ConcreteTypeName{Name: `Void`},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.GenericTypeName{
								Name:    "List",
								Generic: &tree.ConcreteTypeName{Name: `int`},
							},
							Name: &tree.Identifier{Value: `numbers`},
						},
					},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ForEachLoopStatement{
								Field:    &tree.Identifier{Value: `number`},
								Sequence: &tree.Identifier{Value: `numbers`},
								Body: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.ExpressionStatement{
											Expression: &tree.CallExpression{
												Target: &tree.Identifier{Value: `printf`},
												Arguments: tree.CallArgumentList{
													&tree.CallArgument{
														Value: &tree.StringLiteral{Value: "%d "},
													},
													&tree.CallArgument{
														Value: &tree.Identifier{Value: `number`},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Input: `
method add(left int, right int) returns int => left + right
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name: &tree.Identifier{Value: `add`},
					Type: &tree.ConcreteTypeName{Name: `int`},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `int`},
							Name: &tree.Identifier{Value: `left`},
						},
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `int`},
							Name: &tree.Identifier{Value: `right`},
						},
					},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ReturnStatement{
								Value: &tree.BinaryExpression{
									LeftOperand:  &tree.Identifier{Value: `left`},
									RightOperand: &tree.Identifier{Value: `right`},
									Operator:     token.AddOperator,
								},
							},
						},
					},
				},
			},
			{
				Input: `
method greet() => log("Hello")
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name:       &tree.Identifier{Value: `greet`},
					Type:       &tree.ConcreteTypeName{Name: `Void`},
					Parameters: tree.ParameterList{},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ReturnStatement{
								Value: &tree.CallExpression{
									Target: &tree.Identifier{Value: `log`},
									Arguments: tree.CallArgumentList{
										&tree.CallArgument{
											Value: &tree.StringLiteral{Value: `Hello`},
										},
									},
								},
							},
						},
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		})
}

func TestParsing_InvalidMethodDeclaration(testing *testing.T) {
	ExpectError(testing,
		`method call(x int`,
		func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		},
		func(err error) bool {
			return strings.HasSuffix(err.Error(), "expected ) but got: ';'")
		})
	ExpectError(testing,
		`method call(x int,`,
		func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		},
		func(err error) bool {
			return strings.HasSuffix(err.Error(), "expected ) but got: 'eof'")
		})
}
