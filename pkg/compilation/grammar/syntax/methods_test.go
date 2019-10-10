package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/syntax/tree"
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
		method, err := parser.parseMethodDeclaration()
		if err != nil {
			test.Errorf("unexpected error: %s", err)
			continue
		}
		if !compareMethods(*method, expected) {
			test.Errorf("unexpected node value %s, expected %s", *method, expected)
		}
	}
}

func compareMethods(method tree.MethodDeclaration, expected tree.MethodDeclaration) bool {
	if method.Name.Value != expected.Name.Value {
		return false
	}
	if method.Type.FullName() != expected.Type.FullName() {
		return false
	}
	if len(method.Parameters) != len(expected.Parameters) {
		return false
	}
	for index, methodParameter := range method.Parameters {
		expectedParameter := expected.Parameters[index]
		if methodParameter.Name.Value != expectedParameter.Name.Value {
			return false
		}
		if methodParameter.Type.FullName() != expectedParameter.Type.FullName() {
			return false
		}
	}
	return true
}
