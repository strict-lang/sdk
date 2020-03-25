package token

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"testing"
)

func TestPositionDoesNotContainOutside(test *testing.T) {
	entries := map[Position]input.Offset{
		{BeginOffset: 00, EndOffset: 10}: 11,
		{BeginOffset: 11, EndOffset: 12}: 10,
		{BeginOffset: 11, EndOffset: 13}: 14,
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
		{BeginOffset: 0, EndOffset: 10},
		{BeginOffset: 10, EndOffset: 50},
		{BeginOffset: 10, EndOffset: 11},
		{BeginOffset: 10, EndOffset: 12},
	}

	for _, position := range positions {
		for point := position.BeginOffset; point < position.EndOffset; point++ {
			if position.Contains(point) {
				continue
			}
			test.Errorf("position %s does not contain point %d", position, point)
		}
	}
}

func TestPositionContainsItself(test *testing.T) {
	positions := []Position{
		{BeginOffset: 0, EndOffset: 10},
		{BeginOffset: 10, EndOffset: 50},
		{BeginOffset: 10, EndOffset: 11},
		{BeginOffset: 10, EndOffset: 12},
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
		{BeginOffset: 0, EndOffset: 10}: {BeginOffset: 5, EndOffset: 9},
		{BeginOffset: 0, EndOffset: 10}: {BeginOffset: 5, EndOffset: 10},
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
		{BeginOffset: 0, EndOffset: 10}:  {BeginOffset: 5, EndOffset: 11},
		{BeginOffset: 0, EndOffset: 10}:  {BeginOffset: 10, EndOffset: 12},
		{BeginOffset: 10, EndOffset: 15}: {BeginOffset: 5, EndOffset: 10},
	}

	for position, entry := range positions {
		if !position.ContainsPosition(entry) {
			continue
		}
		test.Errorf("position %s does contain %s", position, entry)
	}
}
