package parsing

import (
	scanning2 "gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	"testing"
)

func TestParser_ParseMethodDeclaration(test *testing.T) {
	var entries = map[string]syntaxtree2.MethodDeclaration{
		`
	method number add(number a, number b)	
		ret = a + b
		return ret
`: {
			Name: &syntaxtree2.Identifier{Value: "add"},
			Type: &syntaxtree2.ConcreteTypeName{Name: "number"},
			Parameters: []*syntaxtree2.Parameter{
				{
					Name: &syntaxtree2.Identifier{Value: "a"},
					Type: &syntaxtree2.ConcreteTypeName{Name: "number"},
				},
				{
					Name: &syntaxtree2.Identifier{Value: "b"},
					Type: &syntaxtree2.ConcreteTypeName{Name: "number"},
				},
			},
		},
		`
	method list<number> rangeTo(number)
		for index from 0 to number do
			yield index
`: {
			Name: &syntaxtree2.Identifier{Value: "rangeTo"},
			Type: &syntaxtree2.GenericTypeName{
				Name:    "list",
				Generic: &syntaxtree2.ConcreteTypeName{Name: "number"},
			},
			Parameters: []*syntaxtree2.Parameter{
				{
					Name: &syntaxtree2.Identifier{Value: "number"},
					Type: &syntaxtree2.ConcreteTypeName{Name: "number"},
				},
			},
		},
	}

	for entry, expected := range entries {
		parser := NewTestParser(scanning2.NewStringScanning(entry))
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

func compareMethods(method syntaxtree2.MethodDeclaration, expected syntaxtree2.MethodDeclaration) bool {
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
