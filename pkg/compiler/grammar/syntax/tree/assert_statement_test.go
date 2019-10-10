package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
	"testing"
)

func TestAssertStatement_Accept(testing *testing.T) {
	region := &AssertStatement{
		Region: input.ZeroRegion,
		Expression: &WildcardNode{Region: input.ZeroRegion},
	}
	CreateVisitorTest(region, testing).
		Expect(AssertStatementNodeKind).
		Run()
}

func TestAssertStatement_AcceptRecursive(testing *testing.T) {
	region := &AssertStatement{
		Region:     input.ZeroRegion,
		Expression: &WildcardNode{Region: input.ZeroRegion},
	}
	CreateVisitorTest(region, testing).
		Expect(AssertStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestAssertStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &AssertStatement{
			Region:     input.ZeroRegion,
			Expression: &WildcardNode{Region: input.ZeroRegion},
		}
	})
}