package diagnostic

type Stage struct {
	Name        string
	Alias       string
	Description string
}

var (
	LexicalAnalysis = Stage{
		Name:        "lexical analysis",
		Alias:       "scanning",
		Description: "turns a stream of characters into a stream of tokens",
	}
	SyntacticalAnalysis = Stage{
		Name:        "syntactical analysis",
		Alias:       "grammar",
		Description: "turns a stream of tokens into an abstract grammar tree",
	}
	SemanticAnalysis = Stage{
		Name:        "semantic analysis",
		Description: "walks and validates the abstract grammar tree",
	}
)
