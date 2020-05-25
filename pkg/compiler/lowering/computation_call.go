package lowering

import (
	"github.com/strict-lang/sdk/pkg/compiler/analysis/semantic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

const ComputationCallLoweringPassId = "LetBindingLowering"

func init() {
	passes.Register(newComputationCallLowering())
}

type ComputationCallLowering struct {
	visitor tree.Visitor
}

func newComputationCallLowering() *ComputationCallLowering {
	lowering := &ComputationCallLowering{}
	lowering.visitor = lowering.createVisitor()
	return lowering
}

func (lowering *ComputationCallLowering) Run(context *passes.Context) {
	context.Unit.AcceptRecursive(lowering.visitor)
}

func (lowering *ComputationCallLowering) Id() passes.Id {
	return ComputationCallLoweringPassId
}

func (lowering *ComputationCallLowering) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, semantic.NameResolutionPassId)
}

func (lowering *ComputationCallLowering) createVisitor() tree.Visitor {
	transformer := tree.NewDelegatingExpressionTransformer()
	transformer.CallExpressionVisitor = lowering.transformCall
	return tree.VisitWith(func(node tree.Node) {
		if container, ok := node.(tree.ExpressionContainer); ok {
			container.TransformExpressions(transformer)
		}
	})
}

func (lowering *ComputationCallLowering) transformCall(
	call *tree.CallExpression) tree.Expression {

	if class := lowering.resolveComputationFactoryCallClass(call); class != nil {
		computeCall := &tree.CallExpression{}
		chain := &tree.ChainExpression{
			Expressions: []tree.Expression{call, computeCall},
		}
		computeCall.Parent = chain
		target := &tree.Identifier{
			Value:  "Compute",
			Parent: computeCall,
		}
		computeCall.Target = target
		if computeMethod, ok := scope.LookupNamedMethod(class.Scope, "Compute"); ok {
			target.Bind(computeMethod)
		} else {
			log.Printf("could not find Compute method in Computation: %s", class.QualifiedName)
		}
		return chain
	}
	return call
}

func (lowering *ComputationCallLowering) resolveComputationFactoryCallClass(
	call *tree.CallExpression) *scope.Class {

	if name, ok := call.TargetName(); ok && name.IsBound() {
		if method, isMethod := scope.AsMethodSymbol(name.Binding()); isMethod {
			if method.Factory && isComputation(method.ReturnType) {
				return method.ReturnType
			}
		}
	}
	return nil
}

func isComputation(class *scope.Class) bool {
	computationType := scope.Builtins.Computation
	return class.ActualClass.Is(computationType.ActualClass)
}