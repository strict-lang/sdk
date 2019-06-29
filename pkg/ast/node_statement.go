package ast

type ExpressionStatement struct {
	Expression Expression
}

func (expression *ExpressionStatement) Accept(visitor Visitor) {
	visitor.VisitExpressionStatement(expression)
	expression.expression.Accept(visitor)
}

type MethodCall struct {
	Method    *Member
	Arguments []Expression
}

func (call MethodCall) Type() *Type {
	return call.Method.ValueType
}

func (call *MethodCall) Accept(visitor Visitor) {
	visitor.VisitMethodCall(call)
	for _, argument := range call.Arguments {
		argument.Accept(visitor)
	}
}

type BlockStatement struct {
	Children []Node
}

func (block *BlockStatement) Accept(visitor Visitor) {
	visitor.VisitBlockStatement(block)
	for _, statement := range block.Children {
		statement.Accept(visitor)
	}
}

type ConditionalStatement struct {
	Else      Node
	Body      Node
	Condition Expression
}

func (conditional *ConditionalStatement) Accept(visitor Visitor) {
	visitor.VisitConditionalStatement(conditional)
	conditional.Condition.Accept(visitor)
	conditional.Body.Accept(visitor)
	if conditional.Else != nil {
		conditional.Else.Accept(visitor)
	}
}

type ForLoopStatement struct {
	Body           Node
	Increment      Node
	Termination    Node
	Initialization Node
}

func (loop *ForLoopStatement) Accept(visitor Visitor) {
	visitor.VisitForLoopStatement(loop)
	loop.Increment.Accept(visitor)
	loop.Termination.Accept(visitor)
	loop.Initialization.Accept(visitor)
}

type PreIncrementStatement struct {
}

type PostIncrementStatement struct{}
