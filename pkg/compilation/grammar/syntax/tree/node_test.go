package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
	"testing"
)

var regionTestEntries = []input.Region{
	input.CreateEmptyRegion(10),
	input.ZeroRegion,
	input.CreateRegion(0, 10),
}

func RunNodeRegionTest(testing *testing.T, entryFactory func (region input.Region) Node) {
	for _, entryRegion := range regionTestEntries {
		entry := entryFactory(entryRegion)
		if entry.Region() != entryRegion {
			testing.Errorf("Invalid Region(): Expected %s - got %s",
				entryRegion, entry.Region())
		}
	}
}