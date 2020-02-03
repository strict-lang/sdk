package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func createTestIdentifier() *Identifier {
	return &Identifier{
		Value:  "test",
		Region: input.ZeroRegion,
	}
}
func TestIdentifier_Accept(testing *testing.T) {
	entry := createTestIdentifier()
	CreateVisitorTest(entry, testing).Expect(IdentifierNodeKind).Run()
}

func TestIdentifier_AcceptRecursive(testing *testing.T) {
	entry := createTestIdentifier()
	CreateVisitorTest(entry, testing).Expect(IdentifierNodeKind).RunRecursive()
}

func TestIdentifier_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &Identifier{
			Value:  "test",
			Region: region,
		}
	})
}
