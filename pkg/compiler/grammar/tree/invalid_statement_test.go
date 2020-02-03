package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func TestInvalidStatement_Accept(testing *testing.T) {
	entry := &InvalidStatement{}
	CreateVisitorTest(entry, testing).Expect(InvalidStatementNodeKind).Run()
}

func TestInvalidStatement_AcceptRecursive(testing *testing.T) {
	entry := &InvalidStatement{}
	CreateVisitorTest(entry, testing).Expect(InvalidStatementNodeKind).RunRecursive()
}

func TestInvalidStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &InvalidStatement{Region: region}
	})
}
