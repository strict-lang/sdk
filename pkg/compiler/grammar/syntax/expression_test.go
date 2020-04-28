package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"testing"
)

func TestParseFieldSelectExpression(testing *testing.T) {
	ExpectResult(testing,
		`strict.version`,
		&tree.ChainExpression{
			Expressions: []tree.Expression{
				&tree.Identifier{Value: "strict"},
				&tree.Identifier{Value: "version"},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseExpression()
		})
}

func TestBinaryExpressions(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: "1 + 2",
				ExpectedOutput: &tree.BinaryExpression{
					LeftOperand:  &tree.NumberLiteral{Value: "1"},
					RightOperand: &tree.NumberLiteral{Value: "2"},
					Operator:     token.AddOperator,
				},
			},
			{
				Input: "1 + 2 * 3",
				ExpectedOutput: &tree.BinaryExpression{
					LeftOperand: &tree.NumberLiteral{Value: "1"},
					RightOperand: &tree.BinaryExpression{
						LeftOperand:  &tree.NumberLiteral{Value: "2"},
						RightOperand: &tree.NumberLiteral{Value: "3"},
						Operator:     token.MulOperator,
					},
					Operator: token.AddOperator,
				},
			},
			{
				Input: "1 + 2 * (3 + 4)",
				ExpectedOutput: &tree.BinaryExpression{
					LeftOperand: &tree.NumberLiteral{Value: "1"},
					RightOperand: &tree.BinaryExpression{
						LeftOperand: &tree.NumberLiteral{Value: "2"},
						RightOperand: &tree.BinaryExpression{
							LeftOperand:  &tree.NumberLiteral{Value: "3"},
							RightOperand: &tree.NumberLiteral{Value: "4"},
							Operator:     token.AddOperator,
						},
						Operator: token.MulOperator,
					},
					Operator: token.AddOperator,
				},
			},
			{
				Input: "1 + 2 / 3 + 4",
				ExpectedOutput: &tree.BinaryExpression{
					LeftOperand: &tree.BinaryExpression{
						LeftOperand: &tree.NumberLiteral{Value: "1"},
						RightOperand: &tree.BinaryExpression{
							LeftOperand:  &tree.NumberLiteral{Value: "2"},
							RightOperand: &tree.NumberLiteral{Value: "3"},
							Operator:     token.DivOperator,
						},
						Operator: token.AddOperator,
					},
					RightOperand: &tree.NumberLiteral{Value: "4"},
					Operator:     token.AddOperator,
				},
			},
			{
				Input: "1 + 2 * 3 + 4",
				ExpectedOutput: &tree.BinaryExpression{
					LeftOperand: &tree.BinaryExpression{
						LeftOperand: &tree.NumberLiteral{Value: "1"},
						RightOperand: &tree.BinaryExpression{
							LeftOperand:  &tree.NumberLiteral{Value: "2"},
							RightOperand: &tree.NumberLiteral{Value: "3"},
							Operator:     token.MulOperator,
						},
						Operator: token.AddOperator,
					},
					RightOperand: &tree.NumberLiteral{Value: "4"},
					Operator:     token.AddOperator,
				},
			},
		},
		func(parsing *Parsing) tree.Node {
			return parsing.parseExpression()
		})
}

func TestParsing_ParseListSelectExpression(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `list[0]`,
				ExpectedOutput: &tree.ListSelectExpression{
					Index:  &tree.NumberLiteral{Value: `0`},
					Target: &tree.Identifier{Value: `list`},
				},
			},
			{
				Input: `list[indexes[0]]`,
				ExpectedOutput: &tree.ListSelectExpression{
					Index: &tree.ListSelectExpression{
						Index:  &tree.NumberLiteral{Value: `0`},
						Target: &tree.Identifier{Value: `indexes`},
					},
					Target: &tree.Identifier{Value: `list`},
				},
			},
			{
				Input: `GetList()[0]`,
				ExpectedOutput: &tree.ListSelectExpression{
					Index: &tree.NumberLiteral{Value: `0`},
					Target: &tree.CallExpression{
						Target: &tree.Identifier{Value: `GetList`},
					},
				},
			},
		},
		func(parsing *Parsing) tree.Node {
			return parsing.parseExpression()
		})
}

func TestParsing_ParseCallExpression(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `"text".Length()`,
				ExpectedOutput: &tree.ChainExpression{
					Expressions: []tree.Expression{
						&tree.StringLiteral{Value: "text"},
						&tree.CallExpression{
							Target: &tree.Identifier{Value: "Length"},
						},
					},
				},
			},
			{
				Input: `"text".Length.ToString()`,
				ExpectedOutput: &tree.ChainExpression{
					Expressions: []tree.Expression{
						&tree.StringLiteral{Value: "text"},
						&tree.Identifier{Value: "Length"},
						&tree.CallExpression{
							Target: &tree.Identifier{Value: "ToString"},
						},
					},
				},
			},
			{
				Input: `a.b.c`,
				ExpectedOutput: &tree.ChainExpression{
					Expressions: []tree.Expression{
						&tree.Identifier{Value: "a"},
						&tree.Identifier{Value: "b"},
						&tree.Identifier{Value: "c"},
					},
				},
			},
			{
				Input: `a.b().c`,
				ExpectedOutput: &tree.ChainExpression{
					Expressions: []tree.Expression{
						&tree.Identifier{Value: "a"},
						&tree.CallExpression{
							Target: &tree.Identifier{Value: "b"},
						},
						&tree.Identifier{Value: "c"},
					},
				},
			},
			{
				Input: `Chain().Call()`,
				ExpectedOutput: &tree.ChainExpression{
					Expressions: []tree.Expression{
						&tree.CallExpression{
							Target: &tree.Identifier{Value: "Chain"},
						},
						&tree.CallExpression{
							Target: &tree.Identifier{Value: "Call"},
						},
					},
				},
			},
			{
				Input: `list[0]()`,
				ExpectedOutput: &tree.CallExpression{
					Target: &tree.ListSelectExpression{
						Index:  &tree.NumberLiteral{Value: `0`},
						Target: &tree.Identifier{Value: `list`},
					},
				},
			},
			{
				Input: `call(argument=0)`,
				ExpectedOutput: &tree.CallExpression{
					Target: &tree.Identifier{Value: `call`},
					Arguments: tree.CallArgumentList{
						&tree.CallArgument{
							Label: `argument`,
							Value: &tree.NumberLiteral{Value: `0`},
						},
					},
				},
			},
			{
				Input: `call(argument=0, 1)`,
				ExpectedOutput: &tree.CallExpression{
					Target: &tree.Identifier{Value: `call`},
					Arguments: tree.CallArgumentList{
						&tree.CallArgument{
							Label: `argument`,
							Value: &tree.NumberLiteral{Value: `0`},
						},
						&tree.CallArgument{
							Value: &tree.NumberLiteral{Value: `1`},
						},
					},
				},
			},
		},
		func(parsing *Parsing) tree.Node {
			return parsing.parseExpression()
		})
}
