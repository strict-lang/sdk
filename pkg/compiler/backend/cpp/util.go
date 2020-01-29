package cpp

import "strict.dev/sdk/pkg/compiler/grammar/tree"

func ExtractStatements(nodes []tree.Node) (statements []tree.Statement) {
	for _, node := range nodes {
		if statement, isStatement := node.(tree.Statement); isStatement {
			statements = append(statements, statement)
		}
	}
	return statements
}
