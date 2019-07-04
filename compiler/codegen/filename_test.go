package codegen

import "testing"

var entries = map[string]string{
	"main": "main.generated.cc",
	"test": "test.generated.cc",
	"a_bc": "a_bc.generated.cc",
}

func TestFilenameByUnitName(test *testing.T) {
	for entry, expected := range entries {
		filename := FilenameByUnitName(entry)
		if filename != expected {
			test.Errorf("expected %s but got %s", expected, filename)
		}
	}
}

func TestUnitNameByFilename(test *testing.T) {
	for expected, entry := range entries {
		unitName, err := UnitNameByFilename(entry)
		if err != nil {
			test.Errorf("unexpected error: %s", err)
			continue
		}
		if unitName != expected {
			test.Errorf("expected '%s' but got '%s'", expected, unitName)
		}
	}
}
