package analysis

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/isolate"
	passes "strict.dev/sdk/pkg/compiler/pass"
	"strict.dev/sdk/pkg/compiler/scope"
)

const NameResolutionPassId = "NameResolutionPass"

func init() {
	registerPassInstance(&NameResolutionPass{})
}

type NameResolutionPass struct{}

func (pass *NameResolutionPass) Run(context *passes.Context) {
	visitor := pass.createVisitor()
	context.Unit.AcceptRecursive(visitor)
}

func (pass *NameResolutionPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ScopeCreationPassId, SymbolEnterPassId)
}

func (pass *NameResolutionPass) Id() passes.Id {
	return NameResolutionPassId
}

func (pass *NameResolutionPass) createVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.FieldSelectExpressionVisitor = pass.visitFieldSelect
	visitor.IdentifierVisitor = pass.visitIdentifier
	return visitor
}

func (pass *NameResolutionPass) visitIdentifier(identifier *tree.Identifier) {
	if !identifier.IsBound() && !identifier.IsPartOfDeclaration() {
		pass.bindIdentifier(identifier)
	}
}

func (pass *NameResolutionPass) bindIdentifier(identifier *tree.Identifier) {
	surroundingScope := requireNearestScope(identifier)
	point := identifier.ReferencePoint()
	if entries := surroundingScope.Lookup(point); !entries.IsEmpty() {
		symbol := entries.First().Symbol
		identifier.Bind(symbol)
	} else {
		fmt.Printf("Unknown Symbol %s\n", identifier.Value)
	}
}

func (pass *NameResolutionPass) reportUnknownIdentifier(
	identifier *tree.Identifier) {
}

func (pass *NameResolutionPass) visitFieldSelect(
	expression *tree.FieldSelectExpression) {

	chainResolution := &chainResolution{
		currentScope: requireNearestScope(expression),
		currentNode:  expression,
	}
	chainResolution.bindAll()
}

type chainResolution struct {
	currentScope scope.Scope
	currentNode  tree.Node
}

func (resolution *chainResolution) bindAll() {

}
