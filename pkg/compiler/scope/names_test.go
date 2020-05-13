package scope

import "testing"

func TestFindParentNamespaceName(testing *testing.T) {
	expectParentNamespace(testing, "Strict.Base.String.String", "String")
	expectParentNamespace(testing, "Strict.File.Directory", "File")
	expectParentNamespace(testing, "Strict.Foo", "Strict")
	expectNoParentNamespace(testing, "Root")
}

func expectParentNamespace(testing *testing.T, qualifiedName string, expected string) {
	name, ok := findParentNamespaceName(qualifiedName)
	if !ok || name != expected {
		testing.Errorf("invalid parent namespace: resolved: %s, expected: %s", name, expected)
	}
}

func expectNoParentNamespace(testing *testing.T, qualifiedName string) {
	if name, ok := findParentNamespaceName(qualifiedName); ok {
		testing.Errorf("invalid parent namespace: resolved: %s, expected: none", name)
	}
}
