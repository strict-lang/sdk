package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestRangedLoopStatement_Accept(testing *testing.T) {
	entry := &RangedLoopStatement{
		Region: input.ZeroRegion,
		Field:  &Identifier{Value: "test"},
		Begin:  &WildcardNode{Region: input.ZeroRegion},
		End:    &WildcardNode{Region: input.ZeroRegion},
		Body:   &StatementBlock{Region: input.ZeroRegion},
	}
	CreateVisitorTest(entry, testing).Expect(RangedLoopStatementNodeKind).Run()
}

func TestRangedLoopStatement_AcceptRecursive(testing *testing.T) {
	entry := &RangedLoopStatement{
		Region: input.ZeroRegion,
		Field:  &Identifier{Value: "test"},
		Begin:  &StringLiteral{Value: "begin"},
		End:    &NumberLiteral{Value: "100"},
		Body:   &StatementBlock{Region: input.ZeroRegion},
	}
	CreateVisitorTest(entry, testing).
		Expect(RangedLoopStatementNodeKind).
		Expect(IdentifierNodeKind).
		Expect(StringLiteralNodeKind).
		Expect(NumberLiteralNodeKind).
		Expect(StatementBlockNodeKind).
		RunRecursive()
}

func TestRangedLoopStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &RangedLoopStatement{Region: region}
	})
}
