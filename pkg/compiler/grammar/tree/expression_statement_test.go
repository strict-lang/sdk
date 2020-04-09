package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestExpressionStatement_Accept(testing *testing.T) {
	entry := &ExpressionStatement{Expression: &WildcardNode{}}
	CreateVisitorTest(entry, testing).Expect(ExpressionStatementNodeKind).Run()
}

func TestExpressionStatement_AcceptRecursive(testing *testing.T) {
	entry := &ExpressionStatement{Expression: &WildcardNode{}}
	CreateVisitorTest(entry, testing).
		Expect(ExpressionStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestExpressionStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ExpressionStatement{
			Expression: &WildcardNode{Region: region},
		}
	})
}
