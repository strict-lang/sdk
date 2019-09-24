package linemap

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type LineMap struct {
	lines        []lineEntry
	lineOffsets  []source2.Offset
	recentOffset source2.Offset
	recentLine   source2.LineIndex
}

type lineEntry struct {
	length source2.Offset
	index  source2.LineIndex
	offset source2.Offset
}

func (lines *LineMap) LineAtOffset(offset source2.Offset) source2.LineIndex {
	if offset == lines.recentOffset {
		return lines.recentLine
	}
	line := lines.resolveLineAtOffset(offset)
	lines.recentLine = line
	lines.recentOffset = offset
	return line
}

func (lines *LineMap) resolveLineAtOffset(offset source2.Offset) source2.LineIndex {
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
			return source2.LineIndex(middleIndex + 1)
		}
	}
	return source2.LineIndex(lastIndex) + 1
}

func (lines *LineMap) OffsetAtLine(index source2.LineIndex) source2.Offset {
	lineCount := len(lines.lines)
	if index < 0 || int(index) >= lineCount {
		return source2.Offset(0)
	}
	return source2.Offset(lines.lines[index].offset)
}

func (lines *LineMap) PositionAtOffset(offset source2.Offset) source2.Position {
	lineIndex := lines.LineAtOffset(offset)
	line := lines.LineAtIndex(lineIndex)
	return source2.Position{
		Offset: offset,
		Column: offset - line.Offset,
		Line:   line,
	}
}

func (lines *LineMap) LineAtIndex(lineIndex source2.LineIndex) source2.Line {
	lineIndex -= 1
	lineCount := len(lines.lines)
	if lineIndex < 0 || int(lineIndex) >= lineCount {
		return source2.Line{}
	}
	entry := lines.lines[lineIndex]
	return source2.Line{
		Offset: entry.offset,
		Index:  entry.index + 1,
		Length: entry.length,
	}
}

func (lines *LineMap) LineCount() int {
	return len(lines.lines)
}
