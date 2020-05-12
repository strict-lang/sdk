package sad

import (
	"encoding/json"
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"testing"
)

func TestGeneration(testing *testing.T) {
	generation := newGeneration(&tree.TranslationUnit{
		Name:    "Test.Test",
		Class:   &tree.ClassDeclaration{
			Name:       "Test.Test",
			SuperTypes: []tree.TypeName{
				&tree.ConcreteTypeName{Name: "Test.Super"},
			},
			Children:   []tree.Node{
				&tree.MethodDeclaration{
					Name:       &tree.Identifier{Value: "Run"},
					Parameters: tree.ParameterList{
						&tree.Parameter{
							Type:   &tree.ConcreteTypeName{Name: "App.Options"},
							Name:   &tree.Identifier{Value: "options"},
						},
					},
					Body: &tree.StatementBlock{},
				},
				&tree.FieldDeclaration{
					Name:     &tree.Identifier{Value: "log"},
					TypeName: &tree.ConcreteTypeName{Name: "Strict.Log"},
				},
			},
		},
	})
	descriptor := generation.Generate()
	formatted, _ := json.MarshalIndent(descriptor, "  ", "  ")
	fmt.Printf("generated descriptor: %v", string(formatted))
}
