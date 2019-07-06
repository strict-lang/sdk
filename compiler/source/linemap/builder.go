package linemap

import "github.com/BenjaminNitschke/Strict/compiler/source"

type Builder struct {
	nodes  []*Node
	last   *Node
	index  source.LineIndex
	offset source.Offset
}

func (builder *Builder) Append(length source.Offset) {
	node := &Node{
		line: source.Line{
			Offset: builder.offset,
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
	builder.index++
	builder.offset += length
	// TODO(Merlinosayimwen): Add one for linebreak?
}

func (builder *Builder) NewLinemap() *Linemap {
	var nodes []*Node
	copy(nodes, builder.nodes)
	return &Linemap{
		nodes: nodes,
	}
}

func NewBuilder() *Builder {
	return &Builder{}
}
