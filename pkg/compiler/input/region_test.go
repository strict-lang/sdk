package input

import "testing"

type regionEntry struct {
	begin, end Offset
}

func (entry regionEntry) CreateRegion() Region {
	return CreateRegion(entry.begin, entry.end)
}

func TestRegion_Begin(testing *testing.T) {
	entries := map [regionEntry] Offset{
		{0, 10}: 0,
		{0, 0}: 0,
		{1, 0}: 0,
		{10, 9}: 9,
	}
	for entry, expectedBegin := range entries {
		region := entry.CreateRegion()
		if region.Begin() != expectedBegin {
			testing.Errorf("Unexpected Locate Begin(): got %d - expected %d",
				region.Begin(), expectedBegin)
		}
	}
}

func TestRegion_End(testing *testing.T) {
	entries := map [regionEntry] Offset{
		{0, 10}: 10,
		{0, 0}: 0,
		{1, 0}: 1,
		{10, 9}: 10,
	}
	for entry, expectedEnd := range entries {
		region := entry.CreateRegion()
		if region.End() != expectedEnd {
			testing.Errorf("Unexpected Locate End(): got %d - expected %d",
				region.Begin(), expectedEnd)
		}
	}
}

func TestRegion_ContainsOffset_Inside(testing *testing.T) {
	entries := map [regionEntry] []Offset {
		{0, 10}: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	for entry, containedOffsets := range entries {
		region := entry.CreateRegion()
		testOffsetsInRegion(region, containedOffsets, testing)
	}
}

func testOffsetsInRegion(region Region, offsets []Offset, testing *testing.T) {
	for _, offset := range offsets {
		if !region.ContainsOffset(offset) {
			testing.Errorf("Expected %d to be in Locate %s", offset, region)
		}
	}
}

func TestRegion_ContainsOffset_Outside(testing *testing.T) {
	entries := map [regionEntry] []Offset {
		{0, 10}: {-1, 11},
		{0, 0}: {-1, 0, 1},
	}
	for entry, nonContainedOffsets := range entries {
		region := entry.CreateRegion()
		testOffsetsNotInRegion(region, nonContainedOffsets, testing)
	}
}

func testOffsetsNotInRegion(region Region, offsets []Offset, testing *testing.T) {
	for _, offset := range offsets {
		if region.ContainsOffset(offset) {
			testing.Errorf("Expected %s not to be in Locate %s", offset, region)
		}
	}
}

func TestRegion_ContainsRegion_Inside(testing *testing.T) {
	entries := map [regionEntry] []Region {
		{0, 10}: {{0, 1}, {9, 10}, {0, 10}},
	}
	for entry, containedRegions := range entries {
		region := entry.CreateRegion()
		testRegionsInRegion(region, containedRegions, testing)
	}
}

func testRegionsInRegion(region Region, entries []Region, testing *testing.T) {
	for _, entry := range entries {
		if !region.ContainsRegion(entry) {
			testing.Errorf("Expected %s to be in Locate %s", entry, region)
		}
	}
}

func TestRegion_ContainsRegion_Outside(testing *testing.T) {
	entries := map [regionEntry] []Region {
		{0, 10}: {{-1, 1}, {11, 12}, {-1, 10}, {9, 15}},
	}
	for entry, containedRegions := range entries {
		region := entry.CreateRegion()
		testRegionsNotInRegion(region, containedRegions, testing)
	}
}

func testRegionsNotInRegion(region Region, entries []Region, testing *testing.T) {
	for _, entry := range entries {
		if region.ContainsRegion(entry) {
			testing.Errorf("Expected %s not to be in Locate %s", entry, region)
		}
	}
}

func TestRegion_IsEmpty_Empty(testing *testing.T) {
	entries := []regionEntry {
		{0, 0},
		{10, 10},
	}
	for _, entry := range entries {
		if region := entry.CreateRegion(); !region.IsEmpty() {
			testing.Errorf("Locate %s is empty", region)
		}
	}
}

func TestRegion_IsEmpty_NonEmpty(testing *testing.T) {
	entries := []regionEntry {
		{0, 1},
		{0, 10},
		{10,11},
	}
	for _, entry := range entries {
		if region := entry.CreateRegion(); region.IsEmpty() {
			testing.Errorf("Locate %s is not empty", region)
		}
	}
}

func TestRegion_String(testing *testing.T) {
	entries := map[regionEntry] string {
		{0, 1}: "Locate(begin: 0, end: 1)",
		{0, 0}: "Locate(begin: 0, end: 0)",
	}
	for entry, expectedString := range entries {
		region := entry.CreateRegion()
		if region.String() != expectedString {
			testing.Errorf("Locate has invalid String(): got %s - expected %s",
				region.String(), expectedString)
		}
	}
}