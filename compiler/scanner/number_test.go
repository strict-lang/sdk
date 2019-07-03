package scanner

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/token"
	"testing"
)

func TestGatheringValidNumber(test *testing.T) {
	entries := []string {
		"0",
		"10",
		"0x10",
		"0x0f",
		"0xabc",
		"0xCAFEBABE",
		"0b100100",
		"0.1232123",
		"1232521982",
	}

	for _, entry := range entries {
		scanner := NewStringScanner(entry)
		scanned := scanner.Pull()
		if token.IsInvalidToken(scanned) {
			test.Error("invalid token:")
			scanner.recorder.PrintAllEntries(diagnostic.NewTestPrinter(test))
			continue
		}
		if !token.IsLiteralToken(scanned) {
			test.Errorf("unexpected token %s, exptected literal", scanned)
			continue
		}
		if scanned.Value() != entry {
			test.Errorf("unexpected number '%s', expected '%s'", scanned.Value(), entry)
		}
	}
}
