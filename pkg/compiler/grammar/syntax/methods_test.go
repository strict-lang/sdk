package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
	"testing"
)

func TestParsing_ParseMethodDeclaration(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `
method ListNumbers(begin Number, end Number) returns Number[]
  for number in Range(begin, end)
    yield number
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name: &tree.Identifier{Value: `ListNumbers`},
					Type: &tree.ListTypeName{
						Element: &tree.ConcreteTypeName{Name: "Number"},
					},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `Number`},
							Name: &tree.Identifier{Value: `begin`},
						},
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: `Number`},
							Name: &tree.Identifier{Value: `end`},
						},
					},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ForEachLoopStatement{
								Field: &tree.Identifier{Value: `number`},
								Sequence:   &tree.CallExpression{
									Target:    &tree.Identifier{Value: "Range"},
									Arguments: tree.CallArgumentList{
										&tree.CallArgument{
											Value: &tree.Identifier{Value: "begin"},
										},
										&tree.CallArgument{
											Value: &tree.Identifier{Value: "end"},
										},
									},
									Region:    input.Region{},
									Parent:    nil,
								},
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
method printList(numbers Number[])
  for number in numbers
    printf("%d ", number)
`,
				ExpectedOutput: &tree.MethodDeclaration{
					Name: &tree.Identifier{Value: `printList`},
					Type: &tree.ConcreteTypeName{Name: `Void`},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.ListTypeName{
								Element: &tree.ConcreteTypeName{Name: `Number`},
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
		`method call(x Number`,
		func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		},
		func(err error) bool {
			return strings.HasSuffix(err.Error(), "expected ) but got: ';'")
		})
	ExpectError(testing,
		`method call(x Number,`,
		func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		},
		func(err error) bool {
			return strings.HasSuffix(err.Error(), "expected ) but got: 'eof'")
		})
}
