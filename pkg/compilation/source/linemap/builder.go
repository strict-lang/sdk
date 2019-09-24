package linemap

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type Builder struct {
	index   source2.LineIndex
	offset  source2.Offset
	lines   []lineEntry
	offsets []source2.Offset
}

func (builder *Builder) Append(offset, length source2.Offset) {
	entry := lineEntry{
		offset: offset,
		index:  builder.index,
		length: length,
	}
	builder.offset = offset
	builder.index++
	builder.lines = append(builder.lines, entry)
	builder.offsets = append(builder.offsets, offset)
}

func (builder *Builder) NewLineMap() *LineMap {
	return &LineMap{
		lines:       builder.lines,
		lineOffsets: builder.offsets,
	}
}

func NewBuilder() *Builder {
	return &Builder{
		index: 1,
	}
}
