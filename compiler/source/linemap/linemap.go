package linemap

import "github.com/BenjaminNitschke/Strict/compiler/source"

type Linemap struct {
	nodes []*Node
}

func (lines *Linemap) centerNode() *Node {
	return lines.nodes[len(lines.nodes)/2]
}

func (lines *Linemap) LineAtOffset(offset source.Offset) source.LineIndex {
	center := lines.centerNode()
	node, ok := center.FindNode(offset)
	if !ok {
		return 0
	}
	return node.line.Index
}

func (lines *Linemap) OffsetAtLine(index source.LineIndex) source.Offset {
	if index < 0 || len(lines.nodes) > int(index) {
		return lines.nodes[index].line.Offset
	}
	return 0
}

func (lines *Linemap) PositionAtOffset(offset source.Offset) source.Position {
	center := lines.centerNode()
	node, ok := center.FindNode(offset)
	if !ok {
		return source.Position{Offset: offset}
	}
	line := node.line
	return source.Position{
		Offset: offset,
		Line:   line,
		Column: offset - line.Offset,
	}
}

func (lines *Linemap) PositionAtLine(index source.LineIndex) source.Position {
	if index < 0 || len(lines.nodes) > int(index) {
		line := lines.nodes[index].line
		return source.Position{
			Offset: line.Offset,
			Line:   line,
			Column: 0,
		}
	}
	return source.Position{}
}
