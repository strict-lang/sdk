package ast

type Visitor func (Node) (next Visitor)