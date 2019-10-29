package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestParameter_Accept(testing *testing.T) {
	entry := &Parameter{
		Type:   &ConcreteTypeName{Name: "int"},
		Name:   &Identifier{Value: "sum"},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(ParameterNodeKind).Run()
}

func TestParameter_AcceptRecursive(testing *testing.T) {
	entry := &Parameter{
		Type:   &ConcreteTypeName{Name: "int"},
		Name:   &Identifier{Value: "sum"},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(ParameterNodeKind).
		Expect(IdentifierNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		RunRecursive()
}

func TestParameter_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &Parameter{
			Type:   &ConcreteTypeName{Name: "int"},
			Name:   &Identifier{Value: "sum"},
			Region: region,
		}
	})
}
