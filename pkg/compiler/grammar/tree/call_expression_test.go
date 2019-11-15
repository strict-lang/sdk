package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestCallExpression_Accept(testing *testing.T) {
	entry := &CallExpression{
		Target:    &WildcardNode{Region: input.ZeroRegion},
		Arguments: []*CallArgument{},
		Region:    input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(CallExpressionNodeKind).Run()
}

func TestCallExpression_AcceptRecursive(testing *testing.T) {
	entry := &CallExpression{
		Region: input.ZeroRegion,
		Arguments: []*CallArgument{
			{
				Value:  &WildcardNode{Region: input.ZeroRegion},
				Region: input.ZeroRegion,
			},
		},
		Target: &WildcardNode{Region: input.ZeroRegion},
	}
	CreateVisitorTest(entry, testing).
		Expect(CallExpressionNodeKind).
		Expect(WildcardNodeKind).
		Expect(CallArgumentNodeKind).
		Expect(WildcardNodeKind). // From argument
		RunRecursive()
}

func TestCallExpression_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &CallExpression{
			Target:    &WildcardNode{Region: input.ZeroRegion},
			Arguments: []*CallArgument{},
			Region:    region,
		}
	})
}
