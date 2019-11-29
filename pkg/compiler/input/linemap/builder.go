package linemap

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type Builder struct {
	index   input.LineIndex
	offset  input.Offset
	lines   []lineEntry
	offsets []input.Offset
}

func (builder *Builder) Append(text string, offset, length input.Offset) {
	entry := lineEntry{
		offset:  offset,
		index:   builder.index,
		length:  length,
		content: text,
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
