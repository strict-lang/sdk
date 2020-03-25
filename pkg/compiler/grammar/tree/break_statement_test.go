package tree

import (
	"math/rand"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestBreakStatement_Accept(testing *testing.T) {
	region := &BreakStatement{}
	CreateVisitorTest(region, testing).
		Expect(BreakStatementNodeKind).
		Run()
}

func TestBreakStatement_AcceptRecursive(testing *testing.T) {
	region := &BreakStatement{}
	CreateVisitorTest(region, testing).
		Expect(BreakStatementNodeKind).
		RunRecursive()
}

func TestBreakStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &BreakStatement{Region: region}
	})
}

func TestBreakStatement_Matches(testing *testing.T) {
	CreateMatchTest(testing, &BreakStatement{}).
		Matches(func(random *rand.Rand) Node {
			return &BreakStatement{
				Region: createRandomRegion(random),
			}
		})
}
