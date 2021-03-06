package lowering

import (
	"github.com/strict-lang/sdk/pkg/compiler/analysis/semantic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"log"
)

const LetBindingLoweringPassId = "LetBindingLowering"

func init() {
	passes.Register(newLetBindingLowering())
}

type LetBindingLowering struct {
	visitor tree.Visitor
}

func newLetBindingLowering() *LetBindingLowering {
	lowering := &LetBindingLowering{}
	lowering.visitor = lowering.createVisitor()
	return lowering
}

func (lowering *LetBindingLowering) Run(context *passes.Context) {
	context.Unit.AcceptRecursive(lowering.visitor)
}

func (lowering *LetBindingLowering) Id() passes.Id {
	return LetBindingLoweringPassId
}

func (lowering *LetBindingLowering) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, semantic.NameResolutionPassId)
}

func (lowering *LetBindingLowering) createExpressionTransformer() tree.ExpressionTransformer {
	transformer := tree.NewDelegatingExpressionTransformer()
	transformer.LetBindingVisitor = lowering.rewrite
	return transformer
}

func (lowering *LetBindingLowering) createVisitor() tree.Visitor {
	transformer := lowering.createExpressionTransformer()
	return tree.VisitWith(func(node tree.Node) {
		if container, ok := node.(tree.ExpressionContainer); ok {
			container.TransformExpressions(transformer)
		}
	})
}

func (lowering *LetBindingLowering) rewrite(binding *tree.LetBinding) tree.Expression {
	parent := requireParent(binding)
	if statement, ok := parent.(*tree.ExpressionStatement); ok {
		lowering.rewriteBindingStatement(binding, statement)
	} else {
		return lowering.rewriteBindingInStatement(binding)
	}
	return binding
}

func requireParent(node tree.Node) tree.Node {
	if parent, ok := node.EnclosingNode(); ok {
		return parent
	}
	panic("required a parent")
}

func (lowering *LetBindingLowering) rewriteBindingInStatement(
	binding *tree.LetBinding) tree.Expression {

	if block, parentStatement, ok := findParentStatementInBlock(binding); ok {
		lowering.rewriteInBlock(binding, parentStatement, block)
		return convertNamesToExpression(binding.Names)
	}
	return binding
}

func convertNamesToExpression(names []*tree.Identifier) tree.Expression {
	if len(names) == 0 {
		return &tree.WildcardNode{}
	}
	// TODO: Implement list expressions
	return names[0]
}

func (lowering *LetBindingLowering) rewriteBindingStatement(
	binding *tree.LetBinding, statement tree.Statement) {

	if block, ok := resolveParentBlock(statement); ok {
		lowered := lowering.lower(binding, block)
		block.ReplaceExact(statement, lowered)
	}
}

func (lowering *LetBindingLowering) rewriteInBlock(
	binding *tree.LetBinding,
	parentStatement tree.Statement,
	block *tree.StatementBlock) {

	if index, ok := block.FindIndexOfExact(parentStatement); ok {
		lowering.rewriteInBlockAtIndex(binding, index, block)
	} else {
		log.Fatal("Can not locate child in block")
	}
}

func (lowering *LetBindingLowering) rewriteInBlockAtIndex(
	binding *tree.LetBinding, index int, block *tree.StatementBlock) {

	lowered := lowering.lower(binding, block)
	block.InsertBeforeIndex(index, lowered)
}

// TODO: Support destruction
func (lowering *LetBindingLowering) lower(
	binding *tree.LetBinding, parent tree.Node) *tree.AssignStatement {

	resolvedType, _ := binding.ResolvedType()
	assign := &tree.AssignStatement{
		Value:    binding.Expression,
		Operator: token.AssignOperator,
		Region:   binding.Expression.Locate(),
		Parent:   parent,
	}
	field := &tree.FieldDeclaration{
		Name:     binding.Names[0],
		Region:   binding.Locate(),
		TypeName: tree.ParseTypeName(binding.Locate(), resolvedType.ActualClass),
		Parent:   assign,
		Inferred: true,
	}
	assign.Target = field
	return assign
}

func resolveParentBlock(node tree.Node) (*tree.StatementBlock, bool) {
	if parent, ok := node.EnclosingNode(); ok {
		if block, ok := parent.(*tree.StatementBlock); ok {
			return block, true
		}
	}
	return nil, false
}

func findParentStatementInBlock(node tree.Node) (*tree.StatementBlock, tree.Statement, bool) {
	currentParent, hasParent := node.EnclosingNode()
	for hasParent {
		if parent, ok := currentParent.(tree.Statement); ok {
			if block, inBlock := resolveParentBlock(parent); inBlock {
				return block, parent, true
			}
		}
		currentParent, hasParent = currentParent.EnclosingNode()
	}
	return nil, nil, false
}
