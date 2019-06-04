package token

import (
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"testing"
)

func TestPositionDoesNotContainOutside(test *testing.T) {
	entries := map[Position]source.Offset{
		{Begin: 00, End: 10}: 11,
		{Begin: 11, End: 12}: 10,
		{Begin: 11, End: 13}: 14,
	}

	for position, point := range entries {
		if !position.Contains(point) {
			continue
		}
		test.Errorf("position %s contains point %d", position, point)
	}
}

func TestPositionContains(test *testing.T) {
	positions := []Position{
		{Begin: 0, End: 10},
		{Begin: 10, End: 50},
		{Begin: 10, End: 11},
		{Begin: 10, End: 12},
	}

	for _, position := range positions {
		for point := position.Begin; point < position.End; point++ {
			if position.Contains(point) {
				continue
			}
			test.Errorf("position %s does not contain point %d", position, point)
		}
	}
}

func TestPositionContainsItself(test *testing.T) {
	positions := []Position{
		{Begin: 0, End: 10},
		{Begin: 10, End: 50},
		{Begin: 10, End: 11},
		{Begin: 10, End: 12},
	}

	for _, position := range positions {
		if position.ContainsPosition(position) {
			continue
		}
		test.Errorf("position %s does not contain itself", position)
	}
}

func TestPositionContainsPosition(test *testing.T) {
	positions := map[Position]Position{
		{Begin: 0, End: 10}: {Begin: 5, End: 9},
		{Begin: 0, End: 10}: {Begin: 5, End: 10},
	}

	for position, inner := range positions {
		if position.ContainsPosition(inner) {
			continue
		}
		test.Errorf("position %s does not contain %s", position, inner)
	}
}

func TestPositionDoesNotContainOutsidePosition(test *testing.T) {
	positions := map[Position]Position{
		{Begin: 0, End: 10}:  {Begin: 5, End: 11},
		{Begin: 0, End: 10}:  {Begin: 10, End: 12},
		{Begin: 10, End: 15}: {Begin: 5, End: 10},
	}

	for position, entry := range positions {
		if !position.ContainsPosition(entry) {
			continue
		}
		test.Errorf("position %s does contain %s", position, entry)
	}
}
