package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func createTestConcreteName(name string) *ConcreteTypeName {
	return &ConcreteTypeName{
		Name:         name,
		Region: input.ZeroRegion,
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
