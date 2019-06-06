package linemap

import (
	"github.com/BenjaminNitschke/Strict/pkg/source"
)

func offsetsToNode(offsets ...source.Offset) *Node {
	return arrayToNode(offsetsToArray(offsets...))
}

// offsetsToArray converts an array of line offsets to an array of line
// information. The index of an element in the array becomes the lines
// index. The length of a line is the distance to its next lines offset.
// The arguments last element is not a line offset but the length of the
// source file. The arrays size has to be greater than two, otherwise
// and empty array is returned.
func offsetsToArray(offsets ...source.Offset) []source.Line {
	if len(offsets) < 2 {
		return []source.Line{}
	}
	// Creates the array of converted line information which does not include
	// the arguments last element and thus is one element smaller.
	lines := make([]source.Line, len(offsets)-1)
	lastOffset := offsets[len(offsets)-1]
	// Reverse iterate the array to start at the last line. The lastOffset is
	// set to the offset of the current line. This is done, since the last
	// element in the offsets array is not the offset of a line but the total
	// length of the file. (len(offsets) - 2)) is the index of the last line
	// offset element in the array.
	for index := len(offsets) - 2; index >= 0; index-- {
		offset := offsets[index]
		line := source.Line{
			Index:  source.LineIndex(index),
			Offset: offset,
			Length: lastOffset - offset,
		}
		lastOffset = offset
		lines[index] = line
	}
	return lines
}

// arrayToNode creates a Node from an array of source line. The node is a
// linked list of all nodes in the lines array with the exact order. The
// returned node is nil if the array is empty.
func arrayToNode(lines []source.Line) *Node {
	length := len(lines)
	if length == 0 {
		return nil
	}
	// Head is used to prevent nil checks in the loops body, the last field
	// is initially a pointer to head and then overridden in the loop.
	head := &Node{}
	last := head
	for _, line := range lines {
		last.next = &Node{
			line: line,
			last: last,
		}
		last = last.next
	}
	return head.next
}
