package linemap

import (
	"gitlab.com/strict-lang/sdk/compiler/source"
)

type Linemap struct {
	lines        []lineEntry
	lineOffsets []source.Offset
	recentOffset source.Offset
	recentLine   source.LineIndex
}

type lineEntry struct {
	length source.Offset
	index  source.LineIndex
	offset source.Offset
}

func (lines *Linemap) LineAtOffset(offset source.Offset) source.LineIndex {
	if offset == lines.recentOffset {
		return lines.recentLine
	}
	line := lines.resolveLineAtOffset(offset)
	lines.recentLine = line
	lines.recentOffset = offset
	return line
}

func (lines *Linemap) resolveLineAtOffset(offset source.Offset) source.LineIndex {
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
			return source.LineIndex(middleIndex + 1)
		}
	}
	return source.LineIndex(lastIndex) + 1
}

func (lines *Linemap) OffsetAtLine(index source.LineIndex) source.Offset {
	lineCount := len(lines.lines)
	if index < 0 || int(index) >= lineCount {
		return source.Offset(0)
	}
	return source.Offset(lines.lines[index].offset)
}

func (lines *Linemap) PositionAtOffset(offset source.Offset) source.Position {
	lineIndex := lines.LineAtOffset(offset)
	line := lines.LineAtIndex(lineIndex)
	return source.Position{
		Offset: offset,
		Column: offset - line.Offset,
		Line:   line,
	}
}

func (lines *Linemap) LineAtIndex(lineIndex source.LineIndex) source.Line {
	lineCount := len(lines.lines)
	if lineIndex < 0 || int(lineIndex) >= lineCount {
		return source.Line{}
	}
	entry := lines.lines[lineIndex]
	return source.Line{
		Offset: entry.offset,
		Index:  entry.index,
		Length: entry.length,
	}
}

func (lines *Linemap) LineCount() int {
	return len(lines.lines)
}
