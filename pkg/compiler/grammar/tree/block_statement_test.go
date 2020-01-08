package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestBlockStatement_Accept(testing *testing.T) {
	entry := &StatementBlock{Region: input.ZeroRegion}
	CreateVisitorTest(entry, testing).Expect(BlockStatementNodeKind).Run()
}

func TestBlockStatement_AcceptRecursive(testing *testing.T) {
	entry := &StatementBlock{
		Region:   input.ZeroRegion,
		Children: []Statement{&WildcardNode{}},
	}
	CreateVisitorTest(entry, testing).
		Expect(BlockStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestBlockStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &StatementBlock{
			Children: []Statement{},
			Region:   region,
		}
	})
}
