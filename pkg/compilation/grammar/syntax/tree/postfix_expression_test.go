package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
	"testing"
)

func TestPostfixExpression_Accept(testing *testing.T) {

}

func TestPostfixExpression_AcceptRecursive(testing *testing.T) {

}

func TestPostfixExpression_Region(testing *testing.T) {
	region := input.CreateRegion(0, 10)
	entry := &PostfixExpression{
		Operand:      nil,
		Operator:     0,
		NodePosition: region,
	}
	if entry.Region() != region {
		testing.Errorf("Invalid region: got %s - expected %s", entry.Region(), region)
	}
}