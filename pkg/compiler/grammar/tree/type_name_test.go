package tree

import "testing"

var (
	_ TypeName = &ConcreteTypeName{}
	_ Node     = &ConcreteTypeName{}
	_ TypeName = &GenericTypeName{}
	_ Node     = &GenericTypeName{}
)

func createTestConcreteName(name string) *ConcreteTypeName {
	return &ConcreteTypeName{
		Name:         name,
		NodePosition: ZeroArea{},
	}
}

func TestConcreteTypeName_FullName(test *testing.T) {
	entries := []string{
		"abc", "name", "thisIsTheNameOfAType", "nonGeneric",
	}
	for _, entry := range entries {
		typeName := createTestConcreteName(entry)
		if typeName.FullName() != entry {
			test.Errorf("Unexpected full typename: %s, expected %s", typeName.FullName(), entry)
		}
	}
}
