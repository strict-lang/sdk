package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func createTestPostfixExpression() *PostfixExpression {
	return &PostfixExpression{
		Operand:  &WildcardNode{Region: input.ZeroRegion},
		Operator: token.IncrementOperator,
		Region:   input.ZeroRegion,
	}
}

func TestPostfixExpression_Accept(testing *testing.T) {
	node := createTestPostfixExpression()
	CreateVisitorTest(node, testing).Expect(PostfixExpressionNodeKind).Run()
}

func TestPostfixExpression_AcceptRecursive(testing *testing.T) {
	node := createTestPostfixExpression()
	CreateVisitorTest(node, testing).
		Expect(PostfixExpressionNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestPostfixExpression_Region(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &PostfixExpression{
			Operand:  &WildcardNode{Region: region},
			Operator: token.IncrementOperator,
			Region:   region,
		}
	})
}
