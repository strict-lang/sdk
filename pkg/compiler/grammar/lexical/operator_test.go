package lexical

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
	"testing"
)

func TestGatheringOperator(test *testing.T) {
	entries := map[string]token.Operator{
		"+":   token.AddOperator,
		"++":  token.IncrementOperator,
		"+=":  token.AddAssignOperator,
		"+-":  token.AddOperator,
		"--+": token.DecrementOperator,
		"+=,": token.AddAssignOperator,
		"+,=": token.AddOperator,
		">=":  token.GreaterEqualsOperator,
	}
	for entry, operator := range entries {
		scanner := NewStringScanning(entry)
		scanned, err := scanner.gatherOperator()
		if err != nil {
			test.Errorf("unexpected error: %s", err)
		}
		if scanned != operator {
			test.Errorf("expected operator %s but got %s", operator.String(), operator.String())
		}
	}
}
