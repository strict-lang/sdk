package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func createTestStringLiteral() *StringLiteral {
	return &StringLiteral{
		Value:      "test",
		NodeRegion: input.CreateRegion(0, 5),
	}
}

func TestStringLiteral_Accept(testing *testing.T) {
	entry := createTestStringLiteral()
	CreateVisitorTest(entry, testing).Expect(StringLiteralNodeKind).Run()
}

func TestStringLiteral_AcceptRecursive(testing *testing.T) {
	entry := createTestStringLiteral()
	CreateVisitorTest(entry, testing).Expect(StringLiteralNodeKind).RunRecursive()
}

func TestStringLiteral_Region(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &StringLiteral{
			Value:      "test",
			NodeRegion: region,
		}
	})
}

func TestStringLiteral_ToStringLiteral(testing *testing.T) {
	entry := createTestStringLiteral()
	converted, err := entry.ToStringLiteral()
	if err != nil {
		testing.Errorf("StringLiteral can not be converted to %s", entry.Value)
		return
	}
	if converted.Value != entry.Value {
		testing.Errorf("Invalid ToStringLiteral(): got % - expected %s",
			entry.Value, converted.Value)
	}
}

func TestStringLiteral_ToNumberLiteral_ValidNumber(testing *testing.T) {

}

func TestStringLiteral_ToNumberLiteral_InvalidNumber(testing *testing.T) {

}
