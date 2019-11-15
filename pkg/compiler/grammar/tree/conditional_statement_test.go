package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

var _ Statement = &ConditionalStatement{}

func TestConditionalStatement_Accept(testing *testing.T) {
	entry := &ConditionalStatement{
		Condition:   &WildcardNode{},
		Consequence: &WildcardNode{},
		Region:      input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(ConditionalStatementNodeKind).Run()
}

func TestConditionalStatement_AcceptRecursive_WithAlternative(testing *testing.T) {
	entry := &ConditionalStatement{
		Condition: &UnaryExpression{
			Operator: token.NegateOperator,
			Operand:  &WildcardNode{Region: input.ZeroRegion},
			Region:   input.ZeroRegion,
		},
		Consequence: &WildcardNode{Region: input.ZeroRegion},
		Region:      input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(ConditionalStatementNodeKind).
		Expect(UnaryExpressionNodeKind).
		Expect(WildcardNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestConditionalStatement_AcceptRecursive_WithoutAlternative(testing *testing.T) {
	entry := &ConditionalStatement{
		Condition: &UnaryExpression{
			Operator: token.NegateOperator,
			Operand:  &WildcardNode{Region: input.ZeroRegion},
			Region:   input.ZeroRegion,
		},
		Consequence: &WildcardNode{Region: input.ZeroRegion},
		Alternative: &WildcardNode{Region: input.ZeroRegion},
		Region:      input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(ConditionalStatementNodeKind).
		Expect(UnaryExpressionNodeKind).
		Expect(WildcardNodeKind).
		Expect(WildcardNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestConditionalStatement_HasAlternative_WithAlternative(testing *testing.T) {
	entry := &ConditionalStatement{
		Condition:   &WildcardNode{Region: input.ZeroRegion},
		Alternative: &WildcardNode{Region: input.ZeroRegion},
		Consequence: &WildcardNode{Region: input.ZeroRegion},
		Region:      input.ZeroRegion,
	}
	if !entry.HasAlternative() {
		testing.Error("Expected ConditionalStatement to have an alternative")
	}
}

func TestConditionalStatement_HasAlternative_WithoutAlternative(testing *testing.T) {
	entry := &ConditionalStatement{
		Condition:   &WildcardNode{Region: input.ZeroRegion},
		Consequence: &WildcardNode{Region: input.ZeroRegion},
		Region:      input.ZeroRegion,
	}
	if entry.HasAlternative() {
		testing.Error("Expected ConditionalStatement not to have an alternative")
	}
}

func TestConditionalStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ConditionalStatement{
			Condition:   &WildcardNode{Region: input.ZeroRegion},
			Consequence: &WildcardNode{Region: input.ZeroRegion},
			Region:      region,
		}
	})
}

func TestConditionalStatement_IsModifyingControlFlow(testing *testing.T) {
	entry := &ConditionalStatement{}
	if !entry.IsModifyingControlFlow() {
		testing.Error("Expected ConditionalStatement to modify control flow")
	}
}
