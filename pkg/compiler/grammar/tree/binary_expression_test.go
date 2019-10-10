package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestBinaryExpression_Accept(testing *testing.T) {
	entry := &BinaryExpression{
		LeftOperand:  &WildcardNode{Region: input.ZeroRegion},
		RightOperand: &WildcardNode{Region: input.ZeroRegion},
		Operator:     token.AddOperator,
		Region:       input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(BinaryExpressionNodeKind).Run()
}

func TestBinaryExpression_AcceptRecursive(testing *testing.T) {
	entry := &BinaryExpression{
		LeftOperand:  &StringLiteral{Value: "text"},
		RightOperand: &NumberLiteral{Value: "1234"},
		Operator:     token.AddOperator,
		Region:       input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(BinaryExpressionNodeKind).
		Expect(StringLiteralNodeKind).
		Expect(NumberLiteralNodeKind).
		RunRecursive()
}

func TestBinaryExpression_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &BinaryExpression{
			LeftOperand:  &WildcardNode{Region: input.ZeroRegion},
			RightOperand: &WildcardNode{Region: input.ZeroRegion},
			Operator:     token.AddOperator,
			Region:       input.ZeroRegion,
		}
	})
}
