package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func TestListSelectExpression_Accept(testing *testing.T) {
	region := &ListSelectExpression{
		Target: &WildcardNode{Region: input.ZeroRegion},
		Index:  &WildcardNode{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).Expect(ListSelectExpressionNodeKind).Run()
}

func TestListSelectExpression_AcceptRecursive(testing *testing.T) {
	region := &ListSelectExpression{
		Target: &WildcardNode{Region: input.ZeroRegion},
		Index: &NumberLiteral{
			Value:  "10230",
			Region: input.ZeroRegion,
		},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(region, testing).
		Expect(ListSelectExpressionNodeKind).
		Expect(WildcardNodeKind).
		Expect(NumberLiteralNodeKind).
		RunRecursive()
}

func TestListSelectExpression_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ListSelectExpression{
			Region: region,
			Index:  &WildcardNode{Region: input.ZeroRegion},
			Target: &WildcardNode{Region: input.ZeroRegion},
		}
	})
}
