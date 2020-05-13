package scope

import "strings"

func isTopClass(namespace *Namespace, class *Class) bool {
	if namespaceName, ok := findQualifiedNamespaceName(class.QualifiedName); ok {
		if namespaceName != namespace.QualifiedName {
			return false
		}
		return class.Name() == namespace.Name()
	}
	return false
}

func findLastPartOfQualifiedName(name string) (string, bool) {
	lastDot := strings.LastIndex(name, ".")
	if lastDot != -1 && lastDot + 1 < len(name) - 1{
		return name[lastDot + 1:], true
	}
	return "", false
}

func findQualifiedNamespaceName(name string) (string, bool) {
	lastDot := strings.LastIndex(name, ".")
	if lastDot == -1 {
		return "", false
	}
	return name[:lastDot], true
}

func findParentNamespaceName(name string) (string, bool) {
	lastDot := 0
	for index := len(name) - 1; index >= 0; index-- {
		if name[index] == '.' {
			if lastDot != 0{
				return name[index + 1:lastDot], true
			}
			lastDot = index
		}
	}
	if lastDot != 0 {
		return name[:lastDot], true
	}
	return "", false
}


