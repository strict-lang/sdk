package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestEmptyStatement_Accept(testing *testing.T) {
	entry := &EmptyStatement{}
	CreateVisitorTest(entry, testing).Expect(EmptyStatementNodeKind).Run()
}

func TestEmptyStatement_AcceptRecursive(testing *testing.T) {
	entry := &EmptyStatement{}
	CreateVisitorTest(entry, testing).Expect(EmptyStatementNodeKind).RunRecursive()
}

func TestEmptyStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &EmptyStatement{Region: region}
	})
}
