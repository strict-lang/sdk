package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestNumberLiteral_Accept(testing *testing.T) {
	literal := &NumberLiteral{
		Value:  "1234",
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(literal, testing).Expect(NumberLiteralNodeKind).Run()
}

func TestNumberLiteral_AcceptRecursive(testing *testing.T) {
	literal := &NumberLiteral{
		Value:  "1234",
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(literal, testing).Expect(NumberLiteralNodeKind).RunRecursive()
}

func TestNumberLiteral_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &NumberLiteral{
			Value:  "1234",
			Region: region,
		}
	})
}
