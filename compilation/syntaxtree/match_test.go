package syntaxtree

import "testing"

func TestMatches(test *testing.T) {
	entries := map[Node]Node{
		&CallExpression{
			Method: &Identifier{
				Value:        "foo",
				NodePosition: ZeroPosition{},
			},
			Arguments: []Node{
				&Identifier{
					Value:        "bar",
					NodePosition: ZeroPosition{},
				},
			},
			NodePosition: ZeroPosition{},
		}: &CallExpression{
			Method: &Identifier{
				Value:        "foo",
				NodePosition: ZeroPosition{},
			},
			Arguments: []Node{
				&Identifier{
					Value:        "bar",
					NodePosition: ZeroPosition{},
				},
			},
			NodePosition: ZeroPosition{},
		},
	}

	for left, right := range entries {
		if Matches(left, right) {
			continue
		}
		test.Error("Nodes do not match: ")
		Print(left)
		Print(right)
	}
}

func TestNotMatches(test *testing.T) {
	entries := map[Node]Node{
		&CallExpression{
			Method: &Identifier{
				Value:        "foo",
				NodePosition: ZeroPosition{},
			},
			Arguments: []Node{
				&Identifier{
					Value:        "bar",
					NodePosition: ZeroPosition{},
				},
			},
			NodePosition: ZeroPosition{},
		}: &CallExpression{
			Method: &Identifier{
				Value:        "bar",
				NodePosition: ZeroPosition{},
			},
			Arguments: []Node{
				&Identifier{
					Value:        "foo",
					NodePosition: ZeroPosition{},
				},
			},
			NodePosition: ZeroPosition{},
		},
	}

	for left, right := range entries {
		if !Matches(left, right) {
			continue
		}
		test.Error("Nodes should not match: ")
		Print(left)
		Print(right)
	}
}
