package input

import "fmt"

// Locate is the position of a node in the input code. It may span multiple
// lines or even the whole file. Regions are represented using two offsets and
// thus don't give too many information. This is done because almost every tree
// node has a position field and it would have big memory impacts if positions
// are not small in size. In order to get more information of a nodes position,
// the LineMap from the LineMap package is used. It maps offsets to line data
// and is heavily used in diagnostics. To check whether a node spans multiple
// lines, you have to look up both its begin and end offset in the LineMap.
type Region struct {
	begin Offset
	end   Offset
}

// Begin returns the offset to the nodes begin. If the node is an expression,
// it will return the offset to the expressions first character. The begin is
// never greater than the end offset.
func (region Region) Begin() Offset {
	return region.begin
}

// End returns the offset to the nodes end. If the node is an expression, it
// returns the offset to the expressions last character. The end is not smaller
// than the begin.
func (region Region) End() Offset {
	return region.end
}

// ContainsOffset returns true if the target is within the region.
func (region Region) ContainsOffset(target Offset) bool {
	if region.IsEmpty() {
		return false
	}
	return region.begin <= target && region.end >= target
}

// ContainsRegion returns true if the target region is within the region.
func (region Region) ContainsRegion(target Region) bool {
	return region.ContainsOffset(target.begin) && region.ContainsOffset(target.end)
}

// IsEmpty returns true if the region is empty. A Locate is empty if the begin
// and end are the same. Empty regions can not contain an offset or region.
func (region Region) IsEmpty() bool {
	return region.begin == region.end
}

// String returns a string representation of the Locate.
func (region Region) String() string {
	return fmt.Sprintf("Locate(begin: %d, end: %d)", region.begin, region.end)
}

func orderBySize(left, right Offset) (Offset, Offset) {
	if right < left {
		return right, left
	}
	return left, right
}

// CreateRegion creates a non-empty Locate. It ensures that the begin
// is always prior to the end.
func CreateRegion(begin, end Offset) Region {
	actualBegin, actualEnd := orderBySize(begin, end)
	return Region{
		begin: actualBegin,
		end:   actualEnd,
	}
}

// CreateEmptyRegion creates an empty Locate at the offset.
func CreateEmptyRegion(offset Offset) Region {
	return Region{
		begin: offset,
		end:   offset,
	}
}

// ZeroRegion is an empty Locate that is located at the begin of a file.
var ZeroRegion = CreateEmptyRegion(0)
