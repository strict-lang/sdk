package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestUnaryExpression_Accept(testing *testing.T) {
	entry := &UnaryExpression{
		Operator: token.NegateOperator,
		Operand:  &WildcardNode{Region: input.ZeroRegion},
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(UnaryExpressionNodeKind).Run()
}

func TestUnaryExpression_AcceptRecursive(testing *testing.T) {
	entry := &UnaryExpression{
		Operator: token.NegateOperator,
		Operand:  &WildcardNode{Region: input.ZeroRegion},
		Region:   input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(UnaryExpressionNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestUnaryExpression_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &UnaryExpression{
			Operator: token.NegateOperator,
			Operand:  &WildcardNode{Region: input.ZeroRegion},
			Region:   region,
		}
	})
}
