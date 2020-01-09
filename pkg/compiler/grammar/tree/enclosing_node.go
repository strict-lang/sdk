package tree

// The following code could really use some generics. It should be rewritten once we update
// to go 2 (which will introduce generics to the language).

func SearchEnclosingMethod(node Node) (*MethodDeclaration, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if method, isMethod := currentParent.(*MethodDeclaration); isMethod {
			return method, true
		}
		currentParent, _  = currentParent.EnclosingNode()
	}
	return nil, false
}

func SearchEnclosingClass(node Node) (*ClassDeclaration, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if class, isClass := currentParent.(*ClassDeclaration); isClass {
			return class, true
		}
		currentParent, _  = currentParent.EnclosingNode()
	}
	return nil, false
}

type NodeFilter func(Node) bool

func SearchFirstMatchingEnclosingNode(node Node, filter NodeFilter) (Node, bool) {
	currentParent, _ := node.EnclosingNode()
	for currentParent != nil {
		if filter(currentParent) {
			return node, true
		}
		currentParent, _  = currentParent.EnclosingNode()
	}
	return nil, false
}