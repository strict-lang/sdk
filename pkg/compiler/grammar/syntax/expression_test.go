package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"testing"
)

func TestParseFieldSelectExpression(testing *testing.T) {
	ExpectResult(testing,
		`strict.version`,
		&tree.FieldSelectExpression{
			Target:    &tree.Identifier{Value:  "strict"},
			Selection: &tree.Identifier{Value: "version"},
		}, func(parsing *Parsing) tree.Node {
			return parsing.parseExpression()
		})
}
