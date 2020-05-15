package entering

import (
	"github.com/strict-lang/sdk/pkg/compiler/analysis"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

const ScopeCreationPassId = "ScopeCreationPass"

func init() {
	passes.Register(&ScopeCreationPass{})
}

type ScopeCreationPass struct {
	diagnostics  *diagnostic.Bag
	importScope  scope.Scope
	localIdCount int
}

func (pass *ScopeCreationPass) Run(context *passes.Context) {
	visitor := pass.createScopeResolutionVisitor()
	pass.importScope = analysis.RequireInIsolate(context.Isolate).ImportScope
	context.Unit.AcceptRecursive(visitor)
}

func (pass *ScopeCreationPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ParentAssignPassId, GenericResolutionPassId)
}

func (pass *ScopeCreationPass) Id() passes.Id {
	return ScopeCreationPassId
}

func (pass *ScopeCreationPass) createScopeResolutionVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.BlockStatementVisitor = pass.createBlockStatementScope
	visitor.TranslationUnitVisitor = pass.createTranslationUnitScope
	visitor.ClassDeclarationVisitor = pass.createClassDeclarationScope
	visitor.MethodDeclarationVisitor = pass.createMethodDeclarationScope
	visitor.ConstructorDeclarationVisitor = pass.createConstructorDeclarationScope
	return visitor
}


func (pass *ScopeCreationPass) createBlockStatementScope(block *tree.StatementBlock) {
	surroundingScope := requireNearestScope(block)
	localScope := scope.NewLocalScope(pass.nextLocalIdSuffix(),
		block.Region,
		surroundingScope)
	block.UpdateScope(localScope)
}

func (pass *ScopeCreationPass) createMethodDeclarationScope(method *tree.MethodDeclaration) {
	surroundingScope := requireNearestScope(method)
	localScope := scope.NewLocalScope(
		scope.Id(method.Name.Value),
		method.Region,
		surroundingScope)
	method.UpdateScope(localScope)
}

func (pass *ScopeCreationPass) createTranslationUnitScope(unit *tree.TranslationUnit) {
	id := scope.Id(unit.Name)
	unitScope := scope.NewOuterScopeWithRootId(id, pass.importScope)
	unit.UpdateScope(unitScope)
}

func (pass *ScopeCreationPass) createClassDeclarationScope(
	class *tree.ClassDeclaration) {

	surroundingScope := requireNearestScope(class)
	localScope := scope.NewOuterScope(
		pass.nextLocalIdSuffix(),
		surroundingScope)
	class.UpdateScope(localScope)
}

func (pass *ScopeCreationPass) createConstructorDeclarationScope(
	constructor *tree.ConstructorDeclaration) {

	surroundingScope := requireNearestScope(constructor)
	localScope := scope.NewLocalScope(
		pass.nextLocalIdSuffix(),
		constructor.Region,
		surroundingScope)
	constructor.UpdateScope(localScope)
}

func (pass *ScopeCreationPass) nextLocalIdSuffix() scope.Id {
	pass.localIdCount++
	return scope.Id(pass.localIdCount)
}

func requireNearestScope(node tree.Node) scope.Scope {
	if surroundingScope, ok := tree.ResolveNearestScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding scope does not exist: %v", pretty.Format(node))
	return nil
}

func requireNearestMutableScope(
	node tree.Node) scope.MutableScope {

	if surroundingScope, ok := tree.ResolveNearestMutableScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding scope does not exist: %v", pretty.Format(node))
	return nil
}

func ensureScopeIsMutable(anyScope scope.Scope) scope.MutableScope {
	if mutable, ok := anyScope.(scope.MutableScope); ok {
		return mutable
	}
	log.Fatalf("scope is not mutable: %v", anyScope)
	return nil
}
