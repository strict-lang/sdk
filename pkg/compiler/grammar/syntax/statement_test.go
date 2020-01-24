package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"testing"
)

func TestParsing_ParseConditionalStatement(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `
if 1 < 2
  return 3
`,
				ExpectedOutput: &tree.ConditionalStatement{
					Condition: &tree.BinaryExpression{
						LeftOperand:  &tree.NumberLiteral{Value: `1`},
						RightOperand: &tree.NumberLiteral{Value: `2`},
						Operator:     token.SmallerOperator,
					},
					Consequence: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ReturnStatement{
								Value: &tree.NumberLiteral{Value: `3`},
							},
						},
					},
				},
			},
			{
				Input: `
if 1 < 2
  return 3
else if true
  return 2
else 
  return 1
`,
				ExpectedOutput: &tree.ConditionalStatement{
					Condition: &tree.BinaryExpression{
						LeftOperand:  &tree.NumberLiteral{Value: `1`},
						RightOperand: &tree.NumberLiteral{Value: `2`},
						Operator:     token.SmallerOperator,
					},
					Consequence: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ReturnStatement{
								Value: &tree.NumberLiteral{Value: `3`},
							},
						},
					},
					Alternative: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ConditionalStatement{
								Condition: &tree.Identifier{Value: `true`},
								Consequence: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.ReturnStatement{
											Value: &tree.NumberLiteral{Value: `2`},
										},
									},
								},
								Alternative: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.ReturnStatement{
											Value: &tree.NumberLiteral{Value: `1`},
										},
									},
								},
							},
						},
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseStatement()
		})
}

func TestParsing_ParseReturnStatement(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input:          `return`,
				ExpectedOutput: &tree.ReturnStatement{},
			},
			{
				Input: `return 10`,
				ExpectedOutput: &tree.ReturnStatement{
					Value: &tree.NumberLiteral{Value: `10`},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseStatement()
		})
}

func TestParsing_ParseYieldStatement(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `yield 10`,
				ExpectedOutput: &tree.YieldStatement{
					Value: &tree.NumberLiteral{Value: `10`},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseStatement()
		})
}

func TestParsing_ParseImportStatement(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `import "stdio.h"`,
				ExpectedOutput: &tree.ImportStatement{
					Target: &tree.FileImport{Path: `stdio.h`},
				},
			},
			{
				Input: `import "stdio.h" as io`,
				ExpectedOutput: &tree.ImportStatement{
					Target: &tree.FileImport{
						Path: `stdio.h`,
					},
					Alias: &tree.Identifier{Value: `io`},
				},
			},
			{
				Input: `import Random`,
				ExpectedOutput: &tree.ImportStatement{
					Target: &tree.IdentifierChainImport{
						Chain: []string{`Random`},
					},
				},
			},
			{
				Input: `import Strict.Random`,
				ExpectedOutput: &tree.ImportStatement{
					Target: &tree.IdentifierChainImport{
						Chain: []string{`Strict`, `Random`},
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseImportStatement()
		})
}

func TestParsing_ParseAssignStatement(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `this.name = name`,
				ExpectedOutput: &tree.AssignStatement{
					Target: &tree.FieldSelectExpression{
						Target:    &tree.Identifier{Value: `this`},
						Selection: &tree.Identifier{Value: `name`},
					},
					Value:    &tree.Identifier{Value: `name`},
					Operator: token.AssignOperator,
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseStatement()
		})
}

func TestParsing_ParseIndented(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `
for number in Range(0, 100)
  if IsDivisible(number, 2)
    yield number
  if IsDivisible(number, 4)
    yield number - 1
`,
				ExpectedOutput: &tree.ForEachLoopStatement{
					Sequence: &tree.CallExpression{
						Target: &tree.Identifier{Value: `Range`},
						Arguments: tree.CallArgumentList{
								&tree.CallArgument{
									Value:  &tree.NumberLiteral{Value: "0"},
								},
								&tree.CallArgument{
										Value:  &tree.NumberLiteral{Value: "100"},
								},
						},
					},
					Field: &tree.Identifier{Value: `number`},
					Body: &tree.StatementBlock{
						Children: []tree.Statement{
							&tree.ConditionalStatement{
								Condition:   &tree.CallExpression{
									Target: &tree.Identifier{Value: `IsDivisible`},
									Arguments: tree.CallArgumentList{
										&tree.CallArgument{
											Value: &tree.Identifier{Value: "number"},
										},
										&tree.CallArgument{
											Value: &tree.NumberLiteral{Value: "2"},
										},
									},
								},
								Consequence: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.YieldStatement{
											Value:  &tree.Identifier{Value: `number`},
										},
									},
								},
							},
							&tree.ConditionalStatement{
								Condition:   &tree.CallExpression{
									Target: &tree.Identifier{Value: `IsDivisible`},
									Arguments: tree.CallArgumentList{
										&tree.CallArgument{
											Value: &tree.Identifier{Value: "number"},
										},
										&tree.CallArgument{
											Value: &tree.NumberLiteral{Value: "4"},
										},
									},
								},
								Consequence: &tree.StatementBlock{
									Children: []tree.Statement{
										&tree.YieldStatement{
											Value:  &tree.BinaryExpression{
												LeftOperand: &tree.Identifier{Value: `number`},
												RightOperand: &tree.NumberLiteral{Value: `1`},
												Operator:     token.SubOperator,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseStatement()
		})
}