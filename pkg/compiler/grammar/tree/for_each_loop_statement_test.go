package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestForEachLoopStatement_Accept(testing *testing.T) {
	entry := &ForEachLoopStatement{
		Region:   input.ZeroRegion,
		Body:     &WildcardNode{Region: input.ZeroRegion},
		Sequence: &WildcardNode{Region: input.ZeroRegion},
		Field:    &Identifier{Value: "field"},
	}
	CreateVisitorTest(entry, testing).Expect(ForEachLoopStatementNodeKind).Run()
}

func TestForEachLoopStatement_AcceptRecursive(testing *testing.T) {
	entry := &ForEachLoopStatement{
		Region:   input.ZeroRegion,
		Body:     &WildcardNode{Region: input.ZeroRegion},
		Sequence: &StringLiteral{Value: "strict"},
		Field:    &Identifier{Value: "field"},
	}
	CreateVisitorTest(entry, testing).
		Expect(ForEachLoopStatementNodeKind).
		Expect(IdentifierNodeKind).
		Expect(StringLiteralNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestForEachLoopStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ForEachLoopStatement{Region: region}
	})
}
