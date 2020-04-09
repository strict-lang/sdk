package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
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

const fuzzingNilChance = 0.5
const fuzzTestRepeat = 5

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

func createRandomIdentifier(random *rand.Rand) *Identifier {
	return &Identifier{
		Value: generateRandomIdentifierValue(random, identifierCharset),
	}
}

func createExcludingRandomIdentifier(
	random *rand.Rand, exclude ...string) *Identifier {
	for {
		generated := createRandomIdentifier(random)
		for _, excluded := range exclude {
			if generated.Value != excluded {
				return generated
			}
		}
	}
}

const maxRandomIdentifierLength = 18
const minRandomIdentifierLength = 3
const identifierCharset = "abcdefghijklmnopqrstuvwxyz"

func generateRandomIdentifierValue(random *rand.Rand, charset string) string {
	length := generateRandomIntInRange(
		random, minRandomIdentifierLength, maxRandomIdentifierLength)

	return generateRandomString(random, length, charset)
}

func generateRandomString(random *rand.Rand, length int, charset string) string {
	characters := make([]byte, length)
	for index := 0; index < length; index++ {
		characters[index] = charset[random.Intn(len(charset))]
	}
	return string(characters)
}

func generateRandomIntInRange(random *rand.Rand, begin, end int) int {
	return begin + random.Intn(end-begin)
}

type MatchTest struct {
	base    Node
	testing *testing.T
}

func CreateMatchTest(testing *testing.T, base Node) *MatchTest {
	return &MatchTest{
		base:    base,
		testing: testing,
	}
}

type randomNodeFactory func(*rand.Rand) Node

func (test *MatchTest) Matches(factory randomNodeFactory) *MatchTest {
	random := rand.New(randomSource)
	for count := 0; count < matchTestRepeat; count++ {
		nonMatching := factory(random)
		expectNodesMatch(test.testing, test.base, nonMatching)
	}
	return test
}

func (test *MatchTest) Differs(factory randomNodeFactory) *MatchTest {
	random := rand.New(randomSource)
	for count := 0; count < matchTestRepeat; count++ {
		nonMatching := factory(random)
		expectNodesDontMatch(test.testing, test.base, nonMatching)
	}
	return test
}
