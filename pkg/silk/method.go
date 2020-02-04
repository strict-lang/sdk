package silk

import "strict.dev/sdk/pkg/silk/symbol"

type Method struct {
	Reference  symbol.Reference
	Parameters []Parameter
	Access     Access
	Code       *CodeContainer
}

type Parameter struct {
	Field symbol.Reference
}

// CodeContainer holds the basic blocks of a method. It provides fast iteration
// and fast entry/exist lookup. Every method owns a CodeContainer.
type CodeContainer struct {
	// Blocks is a set of all the blocks that are in the container. The elements
	// do not have to follow any particular order. Unreachable blocks are also
	// contained in this set.
	Blocks []*Block
	// Entry is the first block in a method. It has no predecessors and only one
	// successor. It is usually created implicitly.
	Entry *Block
	// Exit is the last block in a method. It usually contains the ReturnOperation
	// and may have multiple predecessors.
	Exit *Block
}

func (container *CodeContainer) Accept(visitor Visitor) {
	for _, block := range container.Blocks {
		block.Accept(visitor)
	}
}
