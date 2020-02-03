package tree

type ExpressionTransformer interface {
	RewriteIdentifier(*Identifier) Expression
	RewriteStringLiteral(*StringLiteral) Expression
	RewriteNumberLiteral(*NumberLiteral) Expression
	RewriteFieldSelectExpression(*FieldSelectExpression) Expression
	RewriteListSelectExpression(*ListSelectExpression) Expression
	RewriteBinaryExpression(*BinaryExpression) Expression
	RewriteUnaryExpression(*UnaryExpression) Expression
	RewritePostfixExpression(*PostfixExpression) Expression
	RewriteCreateExpression(*CreateExpression) Expression
	RewriteCallArgument(*CallArgument) Expression
	RewriteCallExpression(*CallExpression) Expression
	RewriteLetBinding(*LetBinding) Expression
}

type DelegatingExpressionTransformer struct {
	IdentifierVisitor func(node *Identifier) Expression
	StringLiteralVisitor func(node *StringLiteral) Expression
	NumberLiteralVisitor func(node *NumberLiteral) Expression
	FieldSelectExpressionVisitor func(node *FieldSelectExpression) Expression
	ListSelectExpressionVisitor func(node *ListSelectExpression) Expression
	BinaryExpressionVisitor func(node *BinaryExpression) Expression
	UnaryExpressionVisitor func(node *UnaryExpression) Expression
	PostfixExpressionVisitor func(node *PostfixExpression) Expression
	CreateExpressionVisitor func(node *CreateExpression) Expression
	CallArgumentVisitor func(node *CallArgument) Expression
	CallExpressionVisitor func(node *CallExpression) Expression
	LetBindingVisitor func(node *LetBinding) Expression

}

func NewDelegatingExpressionTransformer() *DelegatingExpressionTransformer {
	return &DelegatingExpressionTransformer{
		IdentifierVisitor: func(node *Identifier) Expression {
			return node
		},
		StringLiteralVisitor: func(node *StringLiteral) Expression {
			return node
		},
		NumberLiteralVisitor: func(node *NumberLiteral) Expression {
			return node
		},
		FieldSelectExpressionVisitor: func(node *FieldSelectExpression) Expression {
			return node
		},
		ListSelectExpressionVisitor: func(node *ListSelectExpression) Expression {
			return node
		},
		BinaryExpressionVisitor: func(node *BinaryExpression) Expression {
			return node
		},
		UnaryExpressionVisitor: func(node *UnaryExpression) Expression {
			return node
		},
		PostfixExpressionVisitor: func(node *PostfixExpression) Expression {
			return node
		},
		CreateExpressionVisitor: func(node *CreateExpression) Expression {
			return node
		},
		CallArgumentVisitor: func(node *CallArgument) Expression {
			return node
		},
		CallExpressionVisitor: func(node *CallExpression) Expression {
			return node
		},
		LetBindingVisitor: func(node *LetBinding) Expression {
			return node
		},

	}
}

func (visitor *DelegatingExpressionTransformer) RewriteIdentifier(node *Identifier) Expression {
	return visitor.IdentifierVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteStringLiteral(node *StringLiteral) Expression {
	return visitor.StringLiteralVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteNumberLiteral(node *NumberLiteral) Expression {
	return visitor.NumberLiteralVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteFieldSelectExpression(node *FieldSelectExpression) Expression {
	return visitor.FieldSelectExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteListSelectExpression(node *ListSelectExpression) Expression {
	return visitor.ListSelectExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteBinaryExpression(node *BinaryExpression) Expression {
	return visitor.BinaryExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteUnaryExpression(node *UnaryExpression) Expression {
	return visitor.UnaryExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewritePostfixExpression(node *PostfixExpression) Expression {
	return visitor.PostfixExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteCreateExpression(node *CreateExpression) Expression {
	return visitor.CreateExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteCallArgument(node *CallArgument) Expression {
	return visitor.CallArgumentVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteCallExpression(node *CallExpression) Expression {
	return visitor.CallExpressionVisitor(node)
}
func (visitor *DelegatingExpressionTransformer) RewriteLetBinding(node *LetBinding) Expression {
	return visitor.LetBindingVisitor(node)
}
