package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestFieldDeclaration_Accept(testing *testing.T) {
	entry := &FieldDeclaration{
		Name: &Identifier{
			Value:  "test",
			Region: input.ZeroRegion,
		},
		TypeName: &ConcreteTypeName{Name: "Type"},
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(FieldDeclarationNodeKind).Run()
}

func TestFieldDeclaration_AcceptRecursive(testing *testing.T) {
	entry := &FieldDeclaration{
		Name: &Identifier{
			Value:  "test",
			Region: input.ZeroRegion,
		},
		TypeName: &ConcreteTypeName{Name: "Class"},
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(FieldDeclarationNodeKind).
		Expect(IdentifierNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		RunRecursive()
}

func TestFieldDeclaration_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &FieldDeclaration{
			Name: &Identifier{
				Value:  "test",
				Region: input.ZeroRegion,
			},
			TypeName: &ConcreteTypeName{Name: "Class"},
			Region:   region,
		}
	})
}
