package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestConstructorDeclaration_Accept(testing *testing.T) {
	region := &ConstructorDeclaration{
		Parameters: ParameterList{},
		Body:      &StatementBlock{Region: input.ZeroRegion},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).Expect(ConstructorDeclarationNodeKind).Run()
}

func TestConstructorDeclaration_AcceptRecursive(testing *testing.T) {
	region := &ConstructorDeclaration{
		Parameters: ParameterList{
			&Parameter{
				Type: &ConcreteTypeName{Name: "Test"},
				Name: &Identifier{
					Value:  "test",
					Region: input.ZeroRegion,
				},
				Region: input.ZeroRegion,
			},
		},
		Body:  &StatementBlock{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).
		Expect(ConstructorDeclarationNodeKind).
		Expect(ParameterNodeKind).
		Expect(IdentifierNodeKind).       // Of Parameter
		Expect(ConcreteTypeNameNodeKind). // Of Parameter
		Expect(StatementBlockNodeKind).
		RunRecursive()
}

func TestConstructorDeclaration_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ConstructorDeclaration{
			Parameters: ParameterList{},
			Body:      &StatementBlock{Region: input.ZeroRegion},
			Region:     region,
		}
	})
}
