package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestConstructorDeclaration_Accept(testing *testing.T) {
	region := &ConstructorDeclaration{
		Parameters: ParameterList{},
		Child:      &WildcardNode{Region: input.ZeroRegion},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).Expect(ConstructorDeclarationNodeKind).Run()
}

func TestConstructorDeclaration_AcceptRecursive(testing *testing.T) {
	region := &ConstructorDeclaration{
		Parameters: ParameterList{
			&Parameter{
				Type: createTestConcreteName("Test"),
				Name: &Identifier{
					Value:  "test",
					Region: input.ZeroRegion,
				},
				Region: input.ZeroRegion,
			},
		},
		Child:  &WildcardNode{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).
		Expect(ConstructorDeclarationNodeKind).
		Expect(ParameterNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestConstructorDeclaration_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ConstructorDeclaration{
			Parameters: ParameterList{},
			Child:      &WildcardNode{Region: input.ZeroRegion},
			Region:     region,
		}
	})
}
