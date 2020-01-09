package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

const ScopeResolutionPassId = "ScopeResolutionPass"

type ScopeResolutionPass struct {
	diagnostics  *diagnostic.Bag
	localIdCount int
}

func (pass *ScopeResolutionPass) Run(context *passes.Context) {
	visitor := pass.createScopeResolutionVisitor()
	context.Unit.AcceptRecursive(visitor)
}

func (pass *ScopeResolutionPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ParentAssignPassId)
}

func (pass *ScopeResolutionPass) createScopeResolutionVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.BlockStatementVisitor = pass.createBlockStatementScope
	visitor.TranslationUnitVisitor = pass.createTranslationUnitScope
	visitor.ClassDeclarationVisitor = pass.createClassDeclarationScope
	visitor.MethodDeclarationVisitor = pass.createMethodDeclarationScope
	visitor.ConstructorDeclarationVisitor = pass.createConstructorDeclarationScope
	return visitor
}

func (*ScopeResolutionPass) createMethodDeclarationScope(method *tree.MethodDeclaration) {
	surroundingScope := requireNearestScope(method)
	localScope := scope.NewLocalScope(
		scope.Id(method.Name.Value),
		method.Region,
		surroundingScope)
	method.UpdateScope(localScope)
}

func (pass *ScopeResolutionPass) createBlockStatementScope(
	block *tree.StatementBlock) {

	surroundingScope := requireNearestScope(block)
	localScope := scope.NewLocalScope(
		pass.nextLocalIdSuffix(),
		block.Region,
		surroundingScope)
	block.UpdateScope(localScope)
}

func (pass *ScopeResolutionPass) createTranslationUnitScope(
	unit *tree.TranslationUnit) {}

func (pass *ScopeResolutionPass) createClassDeclarationScope(
	class *tree.ClassDeclaration) {

	surroundingScope := requireNearestScope(class)
	localScope := scope.NewLocalScope(
		pass.nextLocalIdSuffix(),
		class.Region,
		surroundingScope)
	class.UpdateScope(localScope)
}

func (pass *ScopeResolutionPass) createConstructorDeclarationScope(
	constructor *tree.ConstructorDeclaration) {

	surroundingScope := requireNearestScope(constructor)
	localScope := scope.NewLocalScope(
		pass.nextLocalIdSuffix(),
		constructor.Region,
		surroundingScope)
	constructor.UpdateScope(localScope)
}

func (pass *ScopeResolutionPass) nextLocalIdSuffix() scope.Id {
	pass.localIdCount++
	return scope.Id(pass.localIdCount)
}

func requireNearestScope(node tree.Node) scope.Scope {
	if surroundingScope, ok := tree.ResolveNearestScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding scope does not exist: %v", node)
	return nil
}

func requireNearestMutableScope(node tree.Node) scope.MutableScope {
	if surroundingScope, ok := tree.ResolveNearestMutableScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding mutable scope does not exist: %v", node)
	return nil
}

func ensureScopeIsMutable(anyScope scope.Scope) scope.MutableScope {
	if mutable, ok := anyScope.(scope.MutableScope); ok {
		return mutable
	}
	log.Fatal("scope is not mutable")
	return nil
}
