package tree

// Ensures that all the declaration nodes are implementing the
// Node interface.
var (
	_ Node = &MethodDeclaration{}
	_ Node = &Parameter{}
	_ Node = &FieldDeclaration{}
	_ Node = &ClassDeclaration{}
)
