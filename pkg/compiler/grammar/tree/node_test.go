package tree

import (
	"github.com/google/gofuzz"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"math/rand"
	"testing"
)

var regionTestEntries = []input.Region{
	input.CreateEmptyRegion(10),
	input.ZeroRegion,
	input.CreateRegion(0, 10),
}

type nodeFromRegion = func(input.Region) Node

func RunNodeRegionTest(testing *testing.T, entryFactory nodeFromRegion) {
	for _, entryRegion := range regionTestEntries {
		entry := entryFactory(entryRegion)
		if entry.Locate() != entryRegion {
			testing.Errorf("Invalid Locate(): Expected %s - got %s",
				entryRegion, entry.Locate())
		}
	}
}

const matchTestRepeat = 10
var randomSource = rand.NewSource(0xCAFE)

func RunMatchesTest(testing *testing.T, entryFactory nodeFromRegion) {
	random := rand.New(randomSource)
	for count := 0; count < matchTestRepeat; count++ {
		region := createRandomRegion(random)
		base := entryFactory(region)
		nonMatching := entryFactory(region)
		expectNodesMatch(testing, base, nonMatching)
	}
}

func RunDiffersTest(testing *testing.T, node Node, entryFactory nodeFromRegion) {
	random := rand.New(randomSource)
	for count := 0; count < matchTestRepeat; count++ {
		region := createRandomRegion(random)
		nonMatching := entryFactory(region)
		expectNodesDontMatch(testing, node, nonMatching)
	}
}

const fuzzingNilChance = 0.5
const fuzzTestRepeat = 5

func FuzzDiffersTest(testing *testing.T, node Node, prototype Node) {
	fuzzer := fuzz.New().NilChance(fuzzingNilChance).NumElements(1, 1)
	defer func() {
		err := recover()
		testing.Errorf("Failed fuzzing: %s", err)
	}()
	for count := 0; count < fuzzTestRepeat; count++ {
		fuzzer.Fuzz(&prototype)
		expectNodesDontMatch(testing, node, prototype)
	}
}

func expectNodesMatch(testing *testing.T, base Node, matching Node) {
	if !base.Matches(matching) {
		testing.Error("nodes should match")
	}
}

func expectNodesDontMatch(testing *testing.T, base Node, nonMatching Node) {
	if base.Matches(nonMatching) {
		testing.Error("nodes should not match")
	}
}

func createRandomRegion(random *rand.Rand) input.Region {
	return input.CreateRegion(
		input.Offset(random.Int()),
		input.Offset(random.Int()))
}
