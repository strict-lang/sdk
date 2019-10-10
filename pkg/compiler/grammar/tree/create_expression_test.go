package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func createTestCallExpression() *CallExpression {
	return &CallExpression{
		Target:    &WildcardNode{Region: input.ZeroRegion},
		Arguments: []*CallArgument{},
		Region:    input.ZeroRegion,
	}
}

func TestCreateExpression_Accept(testing *testing.T) {
	entry := &CreateExpression{
		Region: input.ZeroRegion,
		Call:   createTestCallExpression(),
		Type:   createTestConcreteName("String"),
	}
	CreateVisitorTest(entry, testing).Expect(CreateExpressionNodeKind).Run()
}

func TestCreateExpression_AcceptRecursive(testing *testing.T) {
	entry := &CreateExpression{
		Region: input.ZeroRegion,
		Call:   createTestCallExpression(),
		Type:   createTestConcreteName("String"),
	}
	CreateVisitorTest(entry, testing).
		Expect(CreateExpressionNodeKind).
		Expect(ConcreteTypeNameNodeKind).
		Expect(CallExpressionNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestCreateExpression_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &CreateExpression{
			Region: region,
			Call:   createTestCallExpression(),
			Type:   createTestConcreteName("String"),
		}
	})
}