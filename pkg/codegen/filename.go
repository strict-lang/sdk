package codegen

import (
	"fmt"
	"regexp"
)

const (
	FileNameFormat    = "%s.generated.hh"
	FileNameRegexp, _ = regexp.Compile("(\\w+)\\.generated\\.hh")
)

// FilenameByUnitName returns the name of a file that belongs to the passed
// translation unit.
func FilenameByUnitName(unitName string) string {
	return fmt.Sprintf(FileNameFormat, unitName)
}

func UnitNameByFilename(filename string) (string, error) {
	return FileNameRegexp.FindString(filename)
}
