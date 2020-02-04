package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
)

type TestStatement struct {
	MethodName string
	Body       *StatementBlock
	Region     input.Region
	Parent     Node
}

func (test *TestStatement) SetEnclosingNode(target Node) {
	test.Parent = target
}

func (test *TestStatement) EnclosingNode() (Node, bool) {
	return test.Parent, test.Parent != nil
}

func (test *TestStatement) Accept(visitor Visitor) {
	visitor.VisitTestStatement(test)
}

func (test *TestStatement) AcceptRecursive(visitor Visitor) {
	test.Accept(visitor)
	test.Body.AcceptRecursive(visitor)
}

func (test *TestStatement) Locate() input.Region {
	return test.Region
}

func (test *TestStatement) Matches(node Node) bool {
	if target, ok := node.(*TestStatement); ok {
		return test.MethodName == target.MethodName &&
			test.Body.Matches(target.Body)
	}
	return false
}
