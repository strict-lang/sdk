package ast

import "reflect"

// Matches returns whether both nodes match.
func Matches(left, right Node) bool {
	return reflect.DeepEqual(left, right)
}
