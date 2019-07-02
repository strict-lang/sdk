package ast

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
	expression.Expression.Accept(visitor)
}

type MethodCall struct {
	Name      Identifier
	Arguments []Node
}

func (call *MethodCall) Accept(visitor *Visitor) {
	visitor.VisitMethodCall(call)
}

type BlockStatement struct {
	Children []Node
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
}

type ConditionalStatement struct {
	Else      Node
	Body      Node
	Condition Node
}

func (conditional *ConditionalStatement) Accept(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
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

type ForeachLoopStatement struct {
	Body Node
	// Variable is the local field, that has the value of
	// the current element of target.
	Variable Node
	// Target is the collection that is iterated.
	Target Node
}

func (loop *ForeachLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitForeachLoopStatement(loop)
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

type YieldStatement struct {
	// Value is the value that is yielded.
	Value Node
}

func (yield *YieldStatement) Accept(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
}
