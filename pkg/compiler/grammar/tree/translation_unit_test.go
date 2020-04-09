package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func createTestTranslationUnit() *TranslationUnit {
	return &TranslationUnit{
		Name: "Test",
		Imports: []*ImportStatement{
			{
				Target: &IdentifierChainImport{},
				Region: input.ZeroRegion,
			},
		},
		Class: &ClassDeclaration{
			Name:       "Test",
			Parameters: []*ClassParameter{},
			SuperTypes: []TypeName{},
			Children:   []Node{&WildcardNode{Region: input.ZeroRegion}},
			Region:     input.ZeroRegion,
		},
		Region: input.Region{},
	}
}

func TestTranslationUnit_Accept(testing *testing.T) {
	entry := createTestTranslationUnit()
	CreateVisitorTest(entry, testing).Expect(TranslationUnitNodeKind).Run()
}

func TestTranslationUnit_AcceptRecursive(testing *testing.T) {
	entry := createTestTranslationUnit()
	CreateVisitorTest(entry, testing).
		Expect(TranslationUnitNodeKind).
		Expect(ImportStatementNodeKind).
		Expect(ClassDeclarationNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestTranslationUnit_Region(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &TranslationUnit{
			Name:    "test",
			Imports: nil,
			Class:   nil,
			Region:  region,
		}
	})
}

func TestTranslationUnit_ToTypeName(testing *testing.T) {

}
