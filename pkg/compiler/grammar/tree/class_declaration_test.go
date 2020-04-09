package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestClassDeclaration_Accept(testing *testing.T) {
	entry := &ClassDeclaration{
		Name:       "test",
		Parameters: []*ClassParameter{},
		SuperTypes: []TypeName{},
		Children:   []Node{},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(ClassDeclarationNodeKind).Run()
}

func TestClassDeclaration_AcceptRecursive(testing *testing.T) {
	entry := &ClassDeclaration{
		Name:       "test",
		Parameters: []*ClassParameter{},
		SuperTypes: []TypeName{
			&ConcreteTypeName{
				Name:   "test",
				Region: input.ZeroRegion,
			},
		},
		Children: []Node{
			&WildcardNode{Region: input.ZeroRegion},
		},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(ClassDeclarationNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestClassDeclaration_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ClassDeclaration{
			Name:       "test",
			Parameters: []*ClassParameter{},
			SuperTypes: []TypeName{},
			Children:   []Node{},
			Region:     region,
		}
	})
}
