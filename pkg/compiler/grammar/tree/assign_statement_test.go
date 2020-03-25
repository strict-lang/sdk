package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestAssignStatement_Accept(testing *testing.T) {
	entry := &AssignStatement{
		Target:   &WildcardNode{Region: input.ZeroRegion},
		Value:    &WildcardNode{Region: input.ZeroRegion},
		Operator: token.AssignOperator,
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(AssignStatementNodeKind).Run()
}

func TestAssignStatement_AcceptRecursive(testing *testing.T) {
	entry := &AssignStatement{
		Target:   &Identifier{Value: "lhs"},
		Value:    &WildcardNode{Region: input.ZeroRegion},
		Operator: token.AssignOperator,
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(AssignStatementNodeKind).
		Expect(IdentifierNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestAssignStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &AssignStatement{
			Target:   &Identifier{Value: "strict"},
			Value:    &StringLiteral{Value: "Cool"},
			Operator: token.AssignOperator,
			Region:   region,
		}
	})
}
