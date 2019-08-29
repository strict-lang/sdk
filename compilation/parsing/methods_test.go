package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"testing"
)

func TestParser_ParseMethodDeclaration(test *testing.T) {
	var entries = map[string]ast.MethodDeclaration{
		`
	method number add(number a, number b)	
		ret = a + b
		return ret
`: {
			Name: &ast.Identifier{Value: "add"},
			Type: &ast.ConcreteTypeName{Name: "number"},
			Parameters: []*ast.Parameter{
				{
					Name: &ast.Identifier{Value: "a"},
					Type: &ast.ConcreteTypeName{Name: "number"},
				},
				{
					Name: &ast.Identifier{Value: "b"},
					Type: &ast.ConcreteTypeName{Name: "number"},
				},
			},
		},
		`
	method list<number> rangeTo(number)
		for index from 0 to number do
			yield index
`: {
			Name: &ast.Identifier{Value: "rangeTo"},
			Type: &ast.GenericTypeName{
				Name:    "list",
				Generic: &ast.ConcreteTypeName{Name: "number"},
			},
			Parameters: []*ast.Parameter{
				{
					Name: &ast.Identifier{Value: "number"},
					Type: &ast.ConcreteTypeName{Name: "number"},
				},
			},
		},
	}

	for entry, expected := range entries {
		parser := NewTestParser(scanning.NewStringScanning(entry))
		method, err := parser.ParseMethodDeclaration()
		if err != nil {
			test.Errorf("unexpected error: %s", err)
			continue
		}
		if !compareMethods(*method, expected) {
			test.Errorf("unexpected node value %s, expected %s", *method, expected)
		}
		ast.Print(method)
	}
}

func compareMethods(method ast.MethodDeclaration, expected ast.MethodDeclaration) bool {
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
