package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestConcreteTypeName_Accept(testing *testing.T) {
	entry := &ConcreteTypeName{}
	CreateVisitorTest(entry, testing).Expect(ConcreteTypeNameNodeKind).Run()
}

func TestConcreteTypeName_AcceptRecursive(testing *testing.T) {
	entry := &ConcreteTypeName{}
	CreateVisitorTest(entry, testing).Expect(ConcreteTypeNameNodeKind).RunRecursive()
}

func TestConcreteTypeName_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ConcreteTypeName{
			Name:   "type",
			Region: region,
		}
	})
}

func TestConcreteTypeName_NonGenericName(testing *testing.T) {
	entries := []*ConcreteTypeName{
		{Name: "SomethingConcrete"},
		{Name: "ThisIsTheNameOfAType"},
		{Name: "int"},
	}
	for _, entry := range entries {
		if entry.BaseName() != entry.Name {
			testing.Errorf("Entry has invalid BaseName(): expected %s - got %s",
				entry.BaseName(), entry.Name)
		}
	}
}
