package codegen

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	FileNameFormat     = "%s.generated.hh"
	FileNameRegexp     = regexp.MustCompile(`(?m)(\w+)\.generated\.hh`)
	ErrInvalidFilename = errors.New("invalid filename format")
)

// FilenameByUnitName returns the name of a file that belongs to the passed
// translation unit.
func FilenameByUnitName(unitName string) string {
	return fmt.Sprintf(FileNameFormat, unitName)
}

func UnitNameByFilename(filename string) (string, error) {
	for _, match := range FileNameRegexp.FindAllString(filename, -1) {
		return match, nil
	}
	return "", ErrInvalidFilename
}
