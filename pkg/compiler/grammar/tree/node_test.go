package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

var regionTestEntries = []input.Region{
	input.CreateEmptyRegion(10),
	input.ZeroRegion,
	input.CreateRegion(0, 10),
}

func RunNodeRegionTest(testing *testing.T, entryFactory func(region input.Region) Node) {
	for _, entryRegion := range regionTestEntries {
		entry := entryFactory(entryRegion)
		if entry.Locate() != entryRegion {
			testing.Errorf("Invalid Locate(): Expected %s - got %s",
				entryRegion, entry.Locate())
		}
	}
}
