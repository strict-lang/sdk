package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestListTypeName_Accept(testing *testing.T) {
	entry := &ListTypeName{}
	CreateVisitorTest(entry, testing).Expect(ListTypeNameNodeKind).Run()
}

func TestListTypeName_AcceptRecursive(testing *testing.T) {
	entry := &ListTypeName{
		Element: &ConcreteTypeName{Name: "int"},
		Region:  input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(ListTypeNameNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		RunRecursive()
}

func TestListTypeName_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ListTypeName{Region: region}
	})
}

func TestListTypeName_FullName(testing *testing.T) {

}

func TestListTypeName_NonGenericName(testing *testing.T) {

}
