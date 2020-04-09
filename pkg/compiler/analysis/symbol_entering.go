package analysis

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"github.com/strict-lang/sdk/pkg/compiler/typing"
	"log"
)

const SymbolEnterPassId = "SymbolEnterPass"

func init() {
	registerPassInstance(&SymbolEnterPass{})
}

// SymbolEnterPass enters symbols into the scope that they are defined in.
// It also ensures that there are no duplicates.
type SymbolEnterPass struct {
	diagnostics        *diagnostic.Bag
	currentClass       *tree.ClassDeclaration
	currentClassSymbol *scope.Class
	currentUnit        *tree.TranslationUnit
}

func (pass *SymbolEnterPass) Run(context *passes.Context) {
	visitor := pass.createVisitor()
	pass.diagnostics = context.Diagnostic
	context.Unit.AcceptRecursive(visitor)
}

func (pass *SymbolEnterPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ScopeCreationPassId, ImportPassId)
}

func (pass *SymbolEnterPass) Id() passes.Id {
	return SymbolEnterPassId
}

func (pass *SymbolEnterPass) createVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.TranslationUnitVisitor = pass.visitTranslationUnit
	visitor.ClassDeclarationVisitor = pass.visitClassDeclaration
	visitor.MethodDeclarationVisitor = pass.visitMethodDeclaration
	visitor.FieldDeclarationVisitor = pass.visitFieldDeclaration
	visitor.LetBindingVisitor = pass.visitLetBinding
	visitor.ForEachLoopStatementVisitor = pass.visitForEachLoopStatement
	visitor.RangedLoopStatementVisitor = pass.visitRangedLoopStatement
	return visitor
}

func (pass *SymbolEnterPass) visitTranslationUnit(unit *tree.TranslationUnit) {}

func (pass *SymbolEnterPass) visitClassDeclaration(
	declaration *tree.ClassDeclaration) {

	pass.currentClass = declaration
	pass.enterClassDeclaration(declaration)
}

func (pass *SymbolEnterPass) enterClassDeclaration(
	declaration *tree.ClassDeclaration) {

	name := declaration.Name
	surroundingScope := requireNearestMutableScope(declaration)
	if pass.ensureNameDoesNotExist(name, declaration, surroundingScope) {
		symbol := pass.newClassSymbol(declaration)
		pass.currentClassSymbol = symbol
		surroundingScope.Insert(symbol)
	}
}

func (pass *SymbolEnterPass) newClassSymbol(
	declaration *tree.ClassDeclaration) *scope.Class {

	return &scope.Class{
		DeclarationName: declaration.Name,
		Scope:           ensureScopeIsMutable(declaration.Scope()),
		ActualClass:     declaration.NewActualClass(),
	}
}

func (pass *SymbolEnterPass) visitMethodDeclaration(
	declaration *tree.MethodDeclaration) {

	declaration.Name.MarkAsPartOfDeclaration()
	parameterSymbols := pass.enterMethodParameters(declaration)
	if symbol, ok := pass.enterMethodToSurroundingScope(declaration); ok {
		symbol.Parameters = parameterSymbols
	}
}

func (pass *SymbolEnterPass) enterMethodParameters(
	method *tree.MethodDeclaration) []*scope.Field {

	symbols := make([]*scope.Field, len(method.Parameters))
	methodScope := ensureScopeIsMutable(method.Scope())
	for index, parameter := range method.Parameters {
		symbols[index] = pass.enterMethodParameter(parameter, methodScope)
	}
	return symbols
}

func (pass *SymbolEnterPass) enterMethodParameter(
	parameter *tree.Parameter,
	methodScope scope.MutableScope) *scope.Field {

	parameter.Name.MarkAsPartOfDeclaration()
	symbol := pass.newFieldSymbolFromParameter(parameter, methodScope)
	if pass.ensureNameDoesNotExist(parameter.Name.Value, parameter, methodScope) {
		methodScope.Insert(symbol)
		return symbol
	}
	return nil
}

func (pass *SymbolEnterPass) newFieldSymbolFromParameter(
	parameter *tree.Parameter,
	surroundingScope scope.MutableScope) *scope.Field {

	class := pass.requireClass(parameter.Type, surroundingScope)
	return &scope.Field{
		Class:           class,
		DeclarationName: parameter.Name.Value,
		Kind:            scope.ParameterField,
	}
}

func (pass *SymbolEnterPass) enterMethodToSurroundingScope(
	method *tree.MethodDeclaration) (*scope.Method, bool) {

	surroundingScope := requireNearestMutableScope(method)
	if pass.ensureNameDoesNotExist(method.Name.Value, method, surroundingScope) {
		symbol := pass.newMethodSymbol(method, surroundingScope)
		surroundingScope.Insert(symbol)
		return symbol, true
	}
	return nil, false
}

func (pass *SymbolEnterPass) newMethodSymbol(
	method *tree.MethodDeclaration,
	surroundingScope scope.MutableScope) *scope.Method {

	class := pass.requireClass(method.Type, surroundingScope)
	return &scope.Method{
		DeclarationName: method.Name.Value,
		ReturnType:      class,
	}
}

// requireClass is used to resolve classes of certain declarations. While this
// pass is mainly inserting stuff into the scope, it has to get a reference of
// the classes for return types and parameter/field types. If no class with the
// given name exists, an error is reported and a "class replacement" is created
// in the targetScope. This replacement is an empty class with the required name.
// By doing this, the name resolution can continue and more semantic errors may
// be provided to the user, even if the method declaration itself is invalid.
func (pass *SymbolEnterPass) requireClass(
	name tree.TypeName, targetScope scope.MutableScope) *scope.Class {

	returnTypePoint := scope.NewReferencePoint(name.BaseName())
	if class, ok := scope.LookupClass(targetScope, returnTypePoint); ok {
		return class
	}
	pass.reportMissingClass(name)
	return pass.createClassReplacementInScope(name, targetScope)
}

func (pass *SymbolEnterPass) reportMissingClass(name tree.TypeName) {
	log.Printf("Class not found %s\n", name.FullName())
}

func (pass *SymbolEnterPass) createClassReplacementInScope(
	name tree.TypeName,
	targetScope scope.MutableScope) *scope.Class {

	symbol := pass.createClassReplacement(name, targetScope)
	targetScope.Insert(symbol)
	return symbol
}

// Class replacements are created when a certain class is not found, in order
// to keep analysing the code.
func (pass *SymbolEnterPass) createClassReplacement(
	name tree.TypeName, parentScope scope.Scope) *scope.Class {

	return &scope.Class{
		DeclarationName: name.BaseName(),
		Scope:           scope.NewOuterScope(scope.Id(name.BaseName()), parentScope),
		ActualClass:     typing.NewEmptyClass(name.BaseName()),
	}
}

func (pass *SymbolEnterPass) visitForEachLoopStatement(
	loop *tree.ForEachLoopStatement) {

	pass.visitUntypedVariable(loop.Field, loop)
}

func (pass *SymbolEnterPass) visitRangedLoopStatement(
	loop *tree.RangedLoopStatement) {

	pass.visitUntypedVariable(loop.Field, loop)
}

func (pass *SymbolEnterPass) visitFieldDeclaration(
	declaration *tree.FieldDeclaration) {

	declaration.Name.MarkAsPartOfDeclaration()
	surroundingScope := requireNearestMutableScope(declaration)
	name := declaration.Name.Value
	if pass.ensureNameDoesNotExist(name, declaration, surroundingScope) {
		pass.enterFieldDeclaration(declaration, surroundingScope)
	}
}

func (pass *SymbolEnterPass) enterFieldDeclaration(
	declaration *tree.FieldDeclaration, scope scope.MutableScope) {

	if isVariable(declaration) {
		pass.enterVariable(declaration, scope)
	} else {
		pass.enterMemberField(declaration, scope)
	}
}

func (pass *SymbolEnterPass) visitLetBinding(binding *tree.LetBinding) {
	for _, name := range binding.Names {
		pass.visitUntypedVariable(name, binding)
	}
}

func (pass *SymbolEnterPass) visitUntypedVariable(name *tree.Identifier, node tree.Node) {
	name.MarkAsPartOfDeclaration()
	surroundingScope := requireNearestMutableScope(node)
	if pass.ensureNameDoesNotExist(name.Value, node, surroundingScope) {
		surroundingScope.Insert(pass.createUntypedVariable(name.Value))
	}
}

func isVariable(declaration *tree.FieldDeclaration) bool {
	return tree.IsInsideOfMethod(declaration)
}

func (pass *SymbolEnterPass) createMemberField(field *tree.FieldDeclaration) *scope.Field {
	fieldScope := requireNearestMutableScope(field)
	return &scope.Field{
		DeclarationName: field.Name.Value,
		Class:           pass.requireClass(field.TypeName, fieldScope),
		Kind:            scope.MemberField,
		EnclosingClass:  pass.currentClassSymbol,
	}
}

func (pass *SymbolEnterPass) enterMemberField(
	field *tree.FieldDeclaration, scope scope.MutableScope) {

	scope.Insert(pass.createMemberField(field))
}

func (pass *SymbolEnterPass) enterVariable(
	variable *tree.FieldDeclaration, targetScope scope.MutableScope) {

	targetScope.Insert(pass.createUntypedVariable(variable.Name.Value))
}

func (pass *SymbolEnterPass) createUntypedVariable(name string) *scope.Field {
	return &scope.Field{
		DeclarationName: name,
		Kind:            scope.VariableField,
	}
}

func (pass *SymbolEnterPass) ensureNameDoesNotExist(
	name string, node tree.Node, surroundingScope scope.Scope) bool {

	point := scope.NewReferencePoint(name)
	if entries := surroundingScope.Lookup(point); !entries.IsEmpty() {
		existingSymbol := entries.First().Symbol
		pass.reportNameCollision(name, node, existingSymbol)
		return false
	}
	return true
}

func (pass *SymbolEnterPass) reportNameCollision(
	name string,
	node tree.Node,
	existingSymbol scope.Symbol) {

}
