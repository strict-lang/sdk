package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func TestForEachLoopStatement_Accept(testing *testing.T) {
	entry := &ForEachLoopStatement{
		Region:   input.ZeroRegion,
		Body:     &StatementBlock{Region: input.ZeroRegion},
		Sequence: &WildcardNode{Region: input.ZeroRegion},
		Field:    &Identifier{Value: "field"},
	}
	CreateVisitorTest(entry, testing).Expect(ForEachLoopStatementNodeKind).Run()
}

func TestForEachLoopStatement_AcceptRecursive(testing *testing.T) {
	entry := &ForEachLoopStatement{
		Region:   input.ZeroRegion,
		Body:     &StatementBlock{Region: input.ZeroRegion},
		Sequence: &StringLiteral{Value: "strict"},
		Field:    &Identifier{Value: "field"},
	}
	CreateVisitorTest(entry, testing).
		Expect(ForEachLoopStatementNodeKind).
		Expect(IdentifierNodeKind).
		Expect(StringLiteralNodeKind).
		Expect(StatementBlockNodeKind).
		RunRecursive()
}

func TestForEachLoopStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ForEachLoopStatement{Region: region}
	})
}
