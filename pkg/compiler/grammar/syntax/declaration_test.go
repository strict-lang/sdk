package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestDeclarationParsing(testing *testing.T) {
	ExpectAllResults(testing, []ParserTestEntry{
		{
			Input: `
implement Super
`,
			ExpectedOutput: &tree.TranslationUnit{
				Name:      "undefined",
				Namespace: "",
				Class:     &tree.ClassDeclaration{
					Name: "undefined",
					SuperTypes: []tree.TypeName{
						&tree.ConcreteTypeName{
							Name: "Super",
						},
					},
				},
				Region:    input.Region{},
			},
		},
	}, func(parsing *Parsing) tree.Node {
		return parsing.parseTranslationUnit()
	})
}
