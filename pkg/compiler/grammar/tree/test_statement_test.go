package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func TestTestStatement_Accept(testing *testing.T) {
	entry := &TestStatement{
		MethodName: "test",
		Body:       &StatementBlock{},
		Region:     input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(TestStatementNodeKind).Run()
}

func TestTestStatement_AcceptRecursive(testing *testing.T) {
	entry := &TestStatement{
		Body:       &StatementBlock{Region: input.ZeroRegion},
		Region:     input.ZeroRegion,
		MethodName: "test",
	}
	CreateVisitorTest(entry, testing).
		Expect(TestStatementNodeKind).
		Expect(StatementBlockNodeKind).
		RunRecursive()
}

func TestTestStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &TestStatement{
			Body:       &StatementBlock{Region: input.ZeroRegion},
			Region:     region,
			MethodName: "test",
		}
	})
}
