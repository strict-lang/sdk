package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"testing"
)

func TestParser_ParseMethodDeclaration(test *testing.T) {
	var entries = map[string]tree.MethodDeclaration{
		`
	method number add(number a, number b)	
		ret = a + b
		return ret
`: {
			Name: &tree.Identifier{Value: "add"},
			Type: &tree.ConcreteTypeName{Name: "number"},
			Parameters: []*tree.Parameter{
				{
					Name: &tree.Identifier{Value: "a"},
					Type: &tree.ConcreteTypeName{Name: "number"},
				},
				{
					Name: &tree.Identifier{Value: "b"},
					Type: &tree.ConcreteTypeName{Name: "number"},
				},
			},
		},
		`
	method list<number> rangeTo(number)
		for index from 0 to number do
			yield index
`: {
			Name: &tree.Identifier{Value: "rangeTo"},
			Type: &tree.GenericTypeName{
				Name:    "list",
				Generic: &tree.ConcreteTypeName{Name: "number"},
			},
			Parameters: []*tree.Parameter{
				{
					Name: &tree.Identifier{Value: "number"},
					Type: &tree.ConcreteTypeName{Name: "number"},
				},
			},
		},
	}

	for entry, expected := range entries {
		parser := NewTestParser(lexical.NewStringScanning(entry))
		method := parser.parseMethodDeclaration()
		if method.Matches(&expected) {
			test.Errorf(
				"unexpected node value %s, expected %s",
				pretty.Format(method),
				pretty.Format(&expected))
		}
	}
}