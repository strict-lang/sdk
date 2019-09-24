package scanning

import (
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"testing"
)

func TestGatheringOperator(test *testing.T) {
	entries := map[string]token2.Operator{
		"+":   token2.AddOperator,
		"++":  token2.IncrementOperator,
		"+=":  token2.AddAssignOperator,
		"+-":  token2.AddOperator,
		"--+": token2.DecrementOperator,
		"+=,": token2.AddAssignOperator,
		"+,=": token2.AddOperator,
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
