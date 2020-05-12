package sad

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"strings"
	"testing"
)

func TestEncode(testing *testing.T) {
	tree := Generate(&tree.TranslationUnit{
		Name: "Test.Test",
		Class: &tree.ClassDeclaration{
			Name: "Test.Test",
			SuperTypes: []tree.TypeName{
				&tree.ConcreteTypeName{Name: "Test.Super"},
			},
			Children: []tree.Node{
				&tree.MethodDeclaration{
					Name: &tree.Identifier{Value: "Run"},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type: &tree.ConcreteTypeName{Name: "App.Options"},
							Name: &tree.Identifier{Value: "options"},
						},
					},
					Body: &tree.StatementBlock{},
				},
				&tree.MethodDeclaration{
					Name: &tree.Identifier{Value: "RunX"},
					Body: &tree.StatementBlock{},
				},
				&tree.FieldDeclaration{
					Name:     &tree.Identifier{Value: "log"},
					TypeName: &tree.ConcreteTypeName{Name: "Strict.Log"},
				},
			},
		},
	})
	output := Encode(&Tree{Classes: []*Class{tree}})
	const members = "c0.1;m2(3.4)5;m6()5;f7.8;\n"
	const symbols = "Test.Test;Test.Super;Run;options;App.Options;Strict.Base.Void;RunX;log;Strict.Log\n"
	const expected = symbols + members
	if output != expected {
		testing.Errorf("unexpected output: \n%s\n expected: \n%s",
			createBlock(output),
			createBlock(expected))
	}
}

func createBlock(text string) string {
	longestLineLength := findLongestLineLength(text)
	separator := strings.Repeat("-", longestLineLength) + "\n"
	return separator + text + separator
}

func findLongestLineLength(text string) int {
	longest := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > longest {
			longest = len(line)
		}
	}
	return longest
}
