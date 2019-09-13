package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"testing"
)

func TestParser_ParseMethodDeclaration(test *testing.T) {
	var entries = map[string]syntaxtree.MethodDeclaration{
		`
	method number add(number a, number b)	
		ret = a + b
		return ret
`: {
			Name: &syntaxtree.Identifier{Value: "add"},
			Type: &syntaxtree.ConcreteTypeName{Name: "number"},
			Parameters: []*syntaxtree.Parameter{
				{
					Name: &syntaxtree.Identifier{Value: "a"},
					Type: &syntaxtree.ConcreteTypeName{Name: "number"},
				},
				{
					Name: &syntaxtree.Identifier{Value: "b"},
					Type: &syntaxtree.ConcreteTypeName{Name: "number"},
				},
			},
		},
		`
	method list<number> rangeTo(number)
		for index from 0 to number do
			yield index
`: {
			Name: &syntaxtree.Identifier{Value: "rangeTo"},
			Type: &syntaxtree.GenericTypeName{
				Name:    "list",
				Generic: &syntaxtree.ConcreteTypeName{Name: "number"},
			},
			Parameters: []*syntaxtree.Parameter{
				{
					Name: &syntaxtree.Identifier{Value: "number"},
					Type: &syntaxtree.ConcreteTypeName{Name: "number"},
				},
			},
		},
	}

	for entry, expected := range entries {
		parser := NewTestParser(scanning.NewStringScanning(entry))
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

func compareMethods(method syntaxtree.MethodDeclaration, expected syntaxtree.MethodDeclaration) bool {
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
