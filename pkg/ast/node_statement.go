package ast

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
	expression.Expression.Accept(visitor)
}

type MethodCall struct {
	Method    *Member
	Arguments []Node
}

func (call MethodCall) Type() *Type {
	return call.Method.ValueType
}

func (call *MethodCall) Accept(visitor *Visitor) {
	visitor.VisitMethodCall(call)
	for _, argument := range call.Arguments {
		argument.Accept(visitor)
	}
}

type BlockStatement struct {
	Children []Node
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
	for _, statement := range block.Children {
		statement.Accept(visitor)
	}
}

type ConditionalStatement struct {
	Else      Node
	Body      Node
	Condition Node
}

func (conditional *ConditionalStatement) Accept(visitor *Visitor) {
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

func (loop *ForLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitForLoopStatement(loop)
	loop.Increment.Accept(visitor)
	loop.Termination.Accept(visitor)
	loop.Initialization.Accept(visitor)
}

type PreIncrementStatement struct {
	Operand   Node
	Decrement bool
}

type PostIncrementStatement struct {
	Operand   Node
	Decrement bool
}

func (preIncrement *PreIncrementStatement) Accept(visitor Visitor) {
	visitor.VisitPreIncrementStatement(preIncrement)
}

func (postIncrement *PostIncrementStatement) Accept(visitor Visitor) {
	visitor.VisitPostIncrementStatement(postIncrement)
}
