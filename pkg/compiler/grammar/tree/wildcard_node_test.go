package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestWildcardNode_Accept(testing *testing.T) {
	entry := &WildcardNode{Region: input.ZeroRegion}
	CreateVisitorTest(entry, testing).Expect(WildcardNodeKind).Run()
}

func TestWildcardNode_AcceptRecursive(testing *testing.T) {
	entry := &WildcardNode{Region: input.ZeroRegion}
	CreateVisitorTest(entry, testing).Expect(WildcardNodeKind).RunRecursive()
}

func TestWildcardNode_Region(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &WildcardNode{Region: region}
	})
}
