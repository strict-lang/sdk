package linemap

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

type LineMap struct {
	lines        []lineEntry
	lineOffsets  []input.Offset
	recentOffset input.Offset
	recentLine   input.LineIndex
}

type lineEntry struct {
	length input.Offset
	index  input.LineIndex
	offset input.Offset
}

func (lines *LineMap) LineAtOffset(offset input.Offset) input.LineIndex {
	if offset == lines.recentOffset {
		return lines.recentLine
	}
	line := lines.resolveLineAtOffset(offset)
	lines.recentLine = line
	lines.recentOffset = offset
	return line
}

func (lines *LineMap) resolveLineAtOffset(offset input.Offset) input.LineIndex {
	firstIndex := 0
	lastIndex := len(lines.lineOffsets) - 1
	for firstIndex <= lastIndex {
		middleIndex := (firstIndex + lastIndex) >> 1
		middleOffset := lines.lineOffsets[middleIndex]

		if middleOffset < offset {
			firstIndex = middleIndex + 1
		} else if middleOffset > offset {
			lastIndex = middleIndex - 1
		} else {
			return input.LineIndex(middleIndex + 1)
		}
	}
	return input.LineIndex(lastIndex) + 1
}

func (lines *LineMap) OffsetAtLine(index input.LineIndex) input.Offset {
	lineCount := len(lines.lines)
	if index < 0 || int(index) >= lineCount {
		return input.Offset(0)
	}
	return input.Offset(lines.lines[index].offset)
}

func (lines *LineMap) PositionAtOffset(offset input.Offset) input.Position {
	lineIndex := lines.LineAtOffset(offset)
	line := lines.LineAtIndex(lineIndex)
	return input.Position{
		Offset: offset,
		Column: offset - line.Offset,
		Line:   line,
	}
}

func (lines *LineMap) LineAtIndex(lineIndex input.LineIndex) input.Line {
	lineIndex -= 1
	lineCount := len(lines.lines)
	if lineIndex < 0 || int(lineIndex) >= lineCount {
		return input.Line{}
	}
	entry := lines.lines[lineIndex]
	return input.Line{
		Offset: entry.offset,
		Index:  entry.index + 1,
		Length: entry.length,
	}
}

func (lines *LineMap) LineCount() int {
	return len(lines.lines)
}
