package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
	"testing"
)

func createTestPostfixExpression() *PostfixExpression {
	return &PostfixExpression{
		Operand:      &WildcardNode{NodeRegion: input.ZeroRegion},
		Operator:     token.IncrementOperator,
		NodePosition: input.ZeroRegion,
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
			Operand:      &WildcardNode{NodeRegion: region},
			Operator:     token.IncrementOperator,
			NodePosition: region,
		}
	})
}