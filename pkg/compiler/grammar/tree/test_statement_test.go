package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestTestStatement_Accept(testing *testing.T) {
	entry := &TestStatement{
		MethodName: "test",
		Child:      &WildcardNode{},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(TestStatementNodeKind).Run()
}

func TestTestStatement_AcceptRecursive(testing *testing.T) {
	entry := &TestStatement{
		Child:      &WildcardNode{Region: input.ZeroRegion},
		Region:     input.ZeroRegion,
		MethodName: "test",
	}
	CreateVisitorTest(entry, testing).
		Expect(TestStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestTestStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &TestStatement{
			Child:      &WildcardNode{Region: input.ZeroRegion},
			Region:     input.ZeroRegion,
			MethodName: "test",
		}
	})
}