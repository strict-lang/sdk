package syntaxtree

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

// Position is the position of a node in the source code. It may span
// multiple lines or even the whole file. Positions are represented
// using two offsets and thus don't give too many information. This
// is done because almost every AST node has a position field and it
// would have big memory impacts if positions are not small in size.
// In order to get more information of a nodes position, the LineMap
// from the LineMap package is used. It maps offsets to line data and
// is heavily used in diagnostics. To check whether a node spans
// multiple lines, you have to look up both its begin and end offset
// in the LineMap.
type Position interface {
	// Begin returns the offset to the nodes begin. If the node is an
	// expression, it will return the offset to the expressions first
	// character. The begin is never greater than the end offset.
	Begin() source2.Offset
	// End returns the offset to the nodes end. If the node is an
	// expression, i twill return the offset to the expressions last
	// character. The end is never smaller than the begin. When comparing
	// the positions of two nodes, favor the begin offset.
	End() source2.Offset
}

// Positioned is implemented by all nodes that have a specific position
// in the source-code, which matters during semantic-analysis.
type Positioned interface {
	// Position returns the area of source code that is covered by the node.
	// The positions of the nodes children should be inside of its position.
	Position() Position
}

type ZeroPosition struct{}

func (ZeroPosition) Begin() source2.Offset {
	return 0
}

func (ZeroPosition) End() source2.Offset {
	return 0
}
