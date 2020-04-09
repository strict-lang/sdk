package lexical

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"testing"
)

func TestGatheringValidNumber(test *testing.T) {
	entries := []string{
		"0",
		"10",
		"0x10",
		"0x0f",
		"0xabc",
		"0xACE",
		"0b100100",
		"0.1232123",
		"1232521982",
	}

	for _, entry := range entries {
		scanner := NewStringScanning(entry)
		scanned := scanner.Pull()
		if !token.IsLiteralToken(scanned) {
			test.Errorf("unexpected token %s, exptected literal %s", scanned, entry)
			continue
		}
		if scanned.Value() != entry {
			test.Errorf("unexpected number '%s', expected '%s'", scanned.Value(), entry)
		}
	}
}
