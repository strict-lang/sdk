package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestReturnStatement_Accept(testing *testing.T) {
	entry := &ReturnStatement{
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(ReturnStatementNodeKind).Run()
}

func TestReturnStatement_AcceptRecursive_ReturningValue(testing *testing.T) {
	entry := &ReturnStatement{
		Region: input.ZeroRegion,
		Value:  &WildcardNode{},
	}
	CreateVisitorTest(entry, testing).
		Expect(ReturnStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestReturnStatement_AcceptRecursive_NoValue(testing *testing.T) {
	entry := &ReturnStatement{}
	CreateVisitorTest(entry, testing).
		Expect(ReturnStatementNodeKind).
		RunRecursive()
}

func TestReturnStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &ReturnStatement{Region: region}
	})
}
