package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"math/rand"
	"testing"
)

func TestAssertStatement_Accept(testing *testing.T) {
	region := &AssertStatement{
		Region:     input.ZeroRegion,
		Expression: &WildcardNode{Region: input.ZeroRegion},
	}
	CreateVisitorTest(region, testing).
		Expect(AssertStatementNodeKind).
		Run()
}

func TestAssertStatement_AcceptRecursive(testing *testing.T) {
	region := &AssertStatement{
		Region:     input.ZeroRegion,
		Expression: &WildcardNode{Region: input.ZeroRegion},
	}
	CreateVisitorTest(region, testing).
		Expect(AssertStatementNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestAssertStatement_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &AssertStatement{
			Region:     region,
			Expression: &WildcardNode{Region: input.ZeroRegion},
		}
	})
}

func TestAssertStatement_Matches(testing *testing.T) {
	CreateMatchTest(testing, &AssertStatement{
		Expression: &Identifier{Value: "true"},
	}).Matches(func(random *rand.Rand) Node {
		return &AssertStatement{
			Expression: &Identifier{Value: "true"},
			Region:     createRandomRegion(random),
		}
	}).Differs(func(random *rand.Rand) Node {
		return &AssertStatement{
			Expression: createExcludingRandomIdentifier(random, "true"),
			Region:     createRandomRegion(random),
		}
	})
}
