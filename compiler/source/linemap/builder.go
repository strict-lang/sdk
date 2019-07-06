package linemap

import (
	"github.com/BenjaminNitschke/Strict/compiler/source"
)

type Builder struct {
	nodes  []*Node
	last   *Node
	index  source.LineIndex
}

func (builder *Builder) offsetOrLength(offset source.Offset) source.Offset {
	builder.index++
	if builder.last != nil {
		return offset - builder.last.line.Offset
	} else {
		return offset
	}
}

func (builder *Builder) Append(offset source.Offset) {
	length := builder.offsetOrLength(offset)
	node := &Node{
		line: source.Line{
			Offset: offset,
			Index:  builder.index,
			Length: length,
		},
		last: builder.last,
	}
	if builder.last != nil {
		builder.last.next = node
	}
	builder.last = node
	builder.nodes = append(builder.nodes, node)

	// TODO(Merlinosayimwen): Add one for linebreak?
}

func (builder *Builder) NewLinemap() *Linemap {
	return &Linemap{
		nodes: builder.nodes,
	}
}

func NewBuilder() *Builder {
	return &Builder{}
}
