package cpp

import "github.com/strict-lang/sdk/pkg/compiler/grammar/tree"

func isApp(class *tree.ClassDeclaration) bool {
	if class.Trait {
		return false
	}
	for _, superType := range class.SuperTypes {
		if superType.FullName() == "App" {
			return true
		}
	}
	return false
}

func findRunMethod(class *tree.ClassDeclaration) *tree.MethodDeclaration {
	for _, member := range class.Children {
		if method, isMethod := member.(*tree.MethodDeclaration); isMethod && isRunMethod(method) {
			return method
		}
	}
	return nil
}

func isRunMethod(method *tree.MethodDeclaration) bool {
	if method.Name.Value != "Run" {
		return false
	}
	if len(method.Parameters) == 0 {
		return true
	}
	if len(method.Parameters) == 1 {
		firstParameter := method.Parameters[0]
		return firstParameter.Type.FullName() == "Options"
	}
	return false
}
