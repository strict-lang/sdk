package diagnostic

type Stage struct {
	Name        string
	Alias       string
}

var (
	LexicalAnalysis = Stage{
		Name:        "lexical analysis",
		Alias:       "scanning",
	}
	SyntacticalAnalysis = Stage{
		Name:        "syntactical analysis",
		Alias:       "grammar",
	}
	SemanticAnalysis = Stage{
		Name:        "semantic analysis",
		Alias:       "analysis",
	}
)
