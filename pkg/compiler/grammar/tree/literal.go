package tree

type Literal interface {
	Expression
	ToStringLiteral() (*StringLiteral, error)
	ToNumberLiteral() (*NumberLiteral, error)
}
