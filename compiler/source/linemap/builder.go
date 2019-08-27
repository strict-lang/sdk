package linemap

import (
	"gitlab.com/strict-lang/sdk/compiler/source"
)

type Builder struct {
	index   source.LineIndex
	offset  source.Offset
	lines   []lineEntry
	offsets []source.Offset
}

func (builder *Builder) Append(offset, length source.Offset) {
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

func (builder *Builder) NewLinemap() *Linemap {
	return &Linemap{
		lines:       builder.lines,
		lineOffsets: builder.offsets,
	}
}

func NewBuilder() *Builder {
	return &Builder{
		index: 1,
	}
}
