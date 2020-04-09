package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestYieldStatement_Accept(testing *testing.T) {
	entry := &YieldStatement{}
	CreateVisitorTest(entry, testing).Expect(YieldStatementNodeKind).Run()
}

func TestYieldStatement_AcceptRecursive(testing *testing.T) {
	entry := &YieldStatement{
		Region: input.ZeroRegion,
		Value:  &WildcardNode{},
	}
	CreateVisitorTest(entry, testing).
		Expect(YieldStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestYieldStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &YieldStatement{
			Region: region,
			Value:  &WildcardNode{},
		}
	})
}
