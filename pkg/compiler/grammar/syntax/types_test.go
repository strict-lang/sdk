package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"testing"
)

func TestParsing_ParseTypeName(testing *testing.T) {
	ExpectAllResults(testing,
		[]ParserTestEntry{
			{
				Input: `int`,
				ExpectedOutput: &tree.ConcreteTypeName{
					Name: "int",
				},
			},
			{
				Input: `list<int>`,
				ExpectedOutput: &tree.GenericTypeName{
					Name: `list`,
					Generic: &tree.ConcreteTypeName{
						Name: `int`,
					},
				},
			},
			{
				Input: `list<int[]>`,
				ExpectedOutput: &tree.GenericTypeName{
					Name: `list`,
					Generic: &tree.ListTypeName{
						Element: &tree.ConcreteTypeName{
							Name: `int`,
						},
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseTypeName()
		})
}
