package analysis

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
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
		log.Printf("Unknown Symbol %s\n", identifier.Value)
	}
}
func (pass *NameResolutionPass) resolveCallExpression(call *tree.CallExpression) {
	if isResolved(call) {
		return
	}
	if name, ok := call.TargetName(); ok && name.IsBound() {
		searchScope := pass.selectResolutionScope(call)
		if entries := searchScope.Lookup(name.ReferencePoint()); !entries.IsEmpty() {
			if methodSymbol, ok := scope.AsMethodSymbol(entries.First().Symbol); ok {
				name.Bind(methodSymbol)
			}
		}
	}
	// TODO: Report error
}

func (pass *NameResolutionPass) selectResolutionScope(node tree.Expression) scope.Scope {
	if chain, ok := tree.SearchEnclosingChain(node); ok {
		index := findIndexInChain(node.Locate().Begin(), chain)
		formerIndex := index - 1
		if formerIndex >= 0 && formerIndex < len(chain.Expressions) {
			if lastType, ok := chain.Expressions[formerIndex].ResolvedType(); ok {
				return lastType.Scope
			}
			return scope.NewEmptyScope("invalid")
		}
	}
	if localScope, ok := tree.ResolveNearestScope(node); ok {
		return localScope
	}
	return scope.NewEmptyScope("invalid")
}

func findIndexInChain(position input.Offset, chain *tree.ChainExpression) int {
	for index, element := range chain.Expressions {
		if element.Locate().Begin() == position {
			return index
		}
	}
	return 0
}

func (pass *NameResolutionPass) reportUnknownIdentifier(
	identifier *tree.Identifier) {
}

func (pass *NameResolutionPass) visitFieldSelect(
	expression *tree.ChainExpression) {

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
