package entering

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "github.com/strict-lang/sdk/pkg/compiler/pass"
)

func init() {
	passes.Register(newGenericResolutionPass())
}

const GenericResolutionPassId = "GenericResolutionPass"

type GenericResolutionPass struct {
	currentClass        *tree.ClassDeclaration
	typeNameTransformer tree.ExpressionTransformer
	generalVisitor      tree.Visitor
}

func newGenericResolutionPass() *GenericResolutionPass {
	pass := &GenericResolutionPass{}
	generalVisitor := tree.NewEmptyVisitor()
	generalVisitor.ClassDeclarationVisitor = pass.visitTopLevelLetBindings
	generalVisitor.GenericTypeNameVisitor = pass.visitGenericTypeName
	typeNameTransformer := tree.NewDelegatingExpressionTransformer()
	typeNameTransformer.LetBindingVisitor = pass.visitGenericLetBindingInTypeName
	pass.typeNameTransformer = typeNameTransformer
	pass.generalVisitor = generalVisitor
	return pass
}

func (pass *GenericResolutionPass) Id() passes.Id {
	return GenericResolutionPassId
}

func (pass *GenericResolutionPass) Dependencies(isolate *isolate.Isolate) passes.Set {
	return passes.ListInIsolate(isolate, ParentAssignPassId)
}

func (pass *GenericResolutionPass) Run(context *passes.Context) {
	pass.currentClass = context.Unit.Class
	context.Unit.AcceptRecursive(pass.generalVisitor)
}

func (pass *GenericResolutionPass) visitGenericTypeName(typeName *tree.GenericTypeName) {
	for _, argument := range typeName.Arguments {
		argument.Expression = argument.Expression.Transform(pass.typeNameTransformer)
	}
}

func (pass *GenericResolutionPass) visitGenericLetBindingInTypeName(
	binding *tree.LetBinding) tree.Expression {

	if isGenericBinding(binding) {
		pass.addParameters(binding)
		return binding.Names[0]
	}
	return binding
}

func (pass *GenericResolutionPass) visitTopLevelLetBindings(class *tree.ClassDeclaration) {
	var newChildren []tree.Node
	for _, child := range class.Children {
		if binding, isBinding := child.(*tree.LetBinding); isBinding {
			if isTopLevelBinding(binding) && isGenericBinding(binding) {
				pass.addParameters(binding)
				continue
			}
		}
		newChildren = append(newChildren, child)
	}
	class.Children = newChildren
}

func (pass *GenericResolutionPass) addParameters(binding *tree.LetBinding) {
	for _, identifier := range binding.Names {
		parameters := &tree.ClassParameter{
			Name:   identifier.Value,
			Parent: pass.currentClass,
		}
		pass.currentClass.Parameters = append(pass.currentClass.Parameters, parameters)
	}
}

const genericKeyword = "generic"

func isGenericBinding(binding *tree.LetBinding) bool {
	if identifier, isIdentifier := binding.Expression.(*tree.Identifier); isIdentifier {
		return identifier.Value == genericKeyword
	}
	return false
}

func isTopLevelBinding(binding *tree.LetBinding) bool {
	return tree.IsInsideOfMethod(binding)
}
