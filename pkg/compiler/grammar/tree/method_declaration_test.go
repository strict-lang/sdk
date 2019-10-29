package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestMethodDeclaration_Accept(testing *testing.T) {
	entry := &MethodDeclaration{
		Name:       &Identifier{Value: "Test"},
		Type:       &ConcreteTypeName{Name:"void"},
		Parameters: ParameterList{},
		Body:       &WildcardNode{Region: input.ZeroRegion},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(MethodDeclarationNodeKind).Run()
}

func TestMethodDeclaration_AcceptRecursive(testing *testing.T) {
	entry := &MethodDeclaration{
		Name: &Identifier{
			Value:  "Test",
			Region: input.ZeroRegion,
		},
		Type: &ConcreteTypeName{Name: "void"},
		Parameters: ParameterList{
			&Parameter{
				Type: &ConcreteTypeName{Name: "int"},
				Name: &Identifier{
					Value:  "count",
					Region: input.ZeroRegion,
				},
				Region: input.ZeroRegion,
			},
		},
		Body:   &WildcardNode{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(MethodDeclarationNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		Expect(ParameterNodeKind).
		Expect(IdentifierNodeKind). // Of Parameter
		Expect(ConcreteTypeNameNodeKind). // of Parameter
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestMethodDeclaration_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &MethodDeclaration{
			Name:       &Identifier{Value: "Test"},
			Type:       &ConcreteTypeName{Name: "void"},
			Parameters: ParameterList{},
			Body:       &WildcardNode{Region: input.ZeroRegion},
			Region:     region,
		}
	})
}
