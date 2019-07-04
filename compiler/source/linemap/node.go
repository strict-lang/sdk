package linemap

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/source"
	"strings"
)

type Node struct {
	line source.Line
	next *Node
	last *Node
}

func (node *Node) String() string {
	var builder strings.Builder
	current := node
	for {
		builder.WriteString(fmt.Sprintf("Node{%s}, ", current.line))
		current = current.next
		if current == nil {
			break
		}
	}
	return builder.String()
}

func (node *Node) findPreviousNode(offset source.Offset) (*Node, bool) {
	if node.line.Contains(offset) {
		return node, true
	}
	if node.last == nil {
		return nil, false
	}
	return node.last.findPreviousNode(offset)
}

func (node *Node) findNextNode(offset source.Offset) (*Node, bool) {
	if node.line.Contains(offset) {
		return node, true
	}
	if node.next == nil {
		return nil, false
	}
	return node.next.findNextNode(offset)
}

func (node *Node) FindNode(offset source.Offset) (*Node, bool) {
	if node, ok := node.findPreviousNode(offset); ok {
		return node, true
	}
	return node.findNextNode(offset)
}
