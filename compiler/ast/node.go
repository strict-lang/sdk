package ast

import (
	"github.com/BenjaminNitschke/Strict/compiler/scope"
	"github.com/BenjaminNitschke/Strict/compiler/source"
)

// Node is implemented by every node of the ast.
type Node interface {
	// Accept invokes the visitor on all of the nodes children.
	Accept(visitor *Visitor)
}

// Named is implemented by all nodes that have a name.
type Named interface {
	// Name returns the nodes name.
	Name() string
}

// Typed is implemented by all nodes that have a time which is known during
// compilation. While the returned pointer may never be nil, it can point to
// the UnknownType value, indicating that the type is now known.
type Typed interface {
	// Type returns a pointer to the nodes type.
	Type() *Type
}

// Scoped is implemented by all nodes that are only visibile in a certain scope
// and may also only see other nodes within that scope.
type Scoped interface {
	// Scope returns a pointer to the nodes scope.
	Scope() *scope.Scope
}

// Position is the position of a node in the source code. It may span
// multiple lines or even the whole file. Positions are represented
// using two offsets and thus don't give too many information. This
// is done because almost every AST node has a position field and it
// would have big memory impacts if positions are not small in size.
// In order to get more information of a nodes position, the Linemap
// from the linemap package is used. It maps offsets to line data and
// is heavily used in diagnostics. To check whether a node spans
// multiple lines, you have to look up both its begin and end offset
// in the Linemap.
type Position interface {
	// Begin returns the offset to the nodes begin. If the node is an
	// expression, it will return the offset to the expressions first
	// character. The begin is never greater than the end offset.
	Begin() source.Offset
	// End returns the offset to the nodes end. If the node is an
	// expression, i twill return the offset to the expressions last
	// character. The end is never smaller than the begin. When comparing
	// the positions of two nodes, favor the begin offset.
	End() source.Offset
}

// Positioned is implemented by all nodes that have a specific position
// in the source-code, which matters during semantical-analysis.
type Positioned interface {
	// Position returns the area of source code that is covered by the node.
	// The positions of the nodes children should be inside of its position.
	Position() Position
}
