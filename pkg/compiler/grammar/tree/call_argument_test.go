package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"testing"
)

func TestCallArgument_Accept(testing *testing.T) {
	entry := &CallArgument{
		Value:  &WildcardNode{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).Expect(CallArgumentNodeKind).Run()
}

func TestCallArgument_AcceptRecursive(testing *testing.T) {
	entry := &CallArgument{
		Value:  &WildcardNode{Region: input.ZeroRegion},
		Region: input.ZeroRegion,
	}
	CreateVisitorTest(entry, testing).
		Expect(CallArgumentNodeKind).
		Expect(WildcardNodeKind).
		RunRecursive()
}

func TestCallArgument_IsLabeled_Labeled(testing *testing.T) {
	entries := []CallArgument{
		{Label: "label"},
		{Label: "_"},
	}
	for _, entry := range entries {
		if !entry.IsLabeled() {
			testing.Error("CallArgument is not labeled")
		}
	}
}

func TestCallArgument_IsLabeled_Unlabeled(testing *testing.T) {
	entries := []CallArgument{
		{Label: ""},
		{},
	}
	for _, entry := range entries {
		if entry.IsLabeled() {
			testing.Error("CallArgument is not labeled")
		}
	}
}

func TestCallArgument_Locate(testing *testing.T) {
	RunNodeRegionTest(testing, func(region input.Region) Node {
		return &CallArgument{
			Value:  &WildcardNode{Region: input.ZeroRegion},
			Region: region,
		}
	})
}
