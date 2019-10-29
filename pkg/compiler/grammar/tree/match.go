package tree

// Matches returns whether both nodes match.
func Matches(left, right Node) bool {
	if left == nil {
		return right == nil
	}
	return left.Matches(right)
}
