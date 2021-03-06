package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestGenericTypeName_Accept(testing *testing.T) {
	entry := &GenericTypeName{}
	CreateVisitorTest(entry, testing).Expect(GenericTypeNameNodeKind).Run()
}

func TestGenericTypeName_AcceptRecursive(testing *testing.T) {
	entry := &GenericTypeName{
		Name:      "Future",
		Arguments: []*Generic{NewIdentifierGeneric(&Identifier{Value: "String"})},
		Region:    input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(GenericTypeNameNodeKind).
		Expect(IdentifierNodeKind).
		RunRecursive()
}

func TestGenericTypeName_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &GenericTypeName{Region: region}
	})
}

func TestGenericTypeName_FullName(testing *testing.T) {

}

func TestGenericTypeName_NonGenericName(testing *testing.T) {

}
