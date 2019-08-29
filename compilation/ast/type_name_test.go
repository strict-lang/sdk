package ast

import "testing"

var (
	_ TypeName = &ConcreteTypeName{}
	_ Node = &ConcreteTypeName{}
	_ TypeName = &GenericTypeName{}
	_ Node = &GenericTypeName{}
)

func NewTestConcreteTypeName(name string) *ConcreteTypeName {
	return &ConcreteTypeName{
		Name:         name,
		NodePosition: ZeroPosition,
	}
}

func TestConcreteTypeName_FullName(test *testing.T) {
	entries := []string {
		"abc", "name", "thisIsTheNameOfAType", "nonGeneric",
	}
	for _, entry := range entries {
		typeName := NewTestConcreteTypeName(entry)
		if typeName.FullName() != entry {
			test.Errorf("Unexpected full typename: %s, expected %s", typeName.FullName(), entry)
		}
	}
}