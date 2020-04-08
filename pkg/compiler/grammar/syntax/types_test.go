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
					Arguments: []*tree.Generic{
						tree.NewIdentifierGeneric(&tree.Identifier{Value: "int"}),
					},
				},
			},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseTypeName()
		})
}
