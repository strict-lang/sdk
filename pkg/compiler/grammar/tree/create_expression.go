package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type CreateExpression struct {
	Region input.Region
	Call   *CallExpression
	Type   TypeName
}

func (create *CreateExpression) Accept(visitor Visitor) {
	visitor.VisitCreateExpression(create)
}

func (create *CreateExpression) AcceptRecursive(visitor Visitor) {
	create.Accept(visitor)
	create.Type.AcceptRecursive(visitor)
	create.Call.AcceptRecursive(visitor)
}

func (create *CreateExpression) Locate() input.Region {
	return create.Region
}

func (create *CreateExpression) Matches(node Node) bool {
	if target, ok := node.(*CreateExpression); ok {
		return create.Call.Matches(target.Call) &&
			create.Type.Matches(target.Type)
	}
	return false
}