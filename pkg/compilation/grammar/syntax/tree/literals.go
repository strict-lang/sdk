package tree

type StringLiteral struct {
	Value        string
	NodePosition InputRegion
}

func (literal *StringLiteral) Accept(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *StringLiteral) AcceptRecursive(visitor *Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *StringLiteral) Area() InputRegion {
	return literal.NodePosition
}

type NumberLiteral struct {
	Value        string
	NodePosition InputRegion
}

func (literal *NumberLiteral) IsFloat() bool {
	return false
}

func (literal *NumberLiteral) Accept(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}

func (literal *NumberLiteral) AcceptRecursive(visitor *Visitor) {
	visitor.VisitNumberLiteral(literal)
}

func (literal *NumberLiteral) Area() InputRegion {
	return literal.NodePosition
}
