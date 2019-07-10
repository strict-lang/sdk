package linemap

import (
	"fmt"
	pretty "github.com/tonnerre/golang-pretty"
	"gitlab.com/strict-lang/sdk/compiler/source"
)

type Builder struct {
	index source.LineIndex
	offset source.Offset
	lines *[]lineEntry
	offsets *[]source.Offset
}

func (builder *Builder) Append(offset source.Offset) {
	length := offset - builder.offset
	entry := lineEntry{
		offset: offset,
		index:  builder.index,
		length: length,
	}
	fmt.Println(entry)
	builder.offset += offset
	builder.index++
	*builder.lines = append(*builder.lines, entry)
	*builder.offsets = append(*builder.offsets, offset)
	// TODO(Merlinosayimwen): Add one for linebreak?
}

func (builder *Builder) NewLinemap() *Linemap {
	pretty.Println(builder.lines)
	pretty.Println(builder.offsets)
	return &Linemap{
		lines: *builder.lines,
		offsetToLine: *builder.offsets,
	}
}

func NewBuilder() *Builder {
	return &Builder{
		index: 1,
		lines: &[]lineEntry{},
		offsets: &[]source.Offset{},
	}
}
