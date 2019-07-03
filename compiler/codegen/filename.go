package codegen

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	FileNameFormat     = "%s.generated.cc"
	FileNameRegexp     = regexp.MustCompile(`(?P<Unit>\w+)\.generated\.cc`)
	ErrInvalidFilename = errors.New("invalid filename format")
)

// FilenameByUnitName returns the name of a file that belongs to the passed
// translation unit.
func FilenameByUnitName(unitName string) string {
	return fmt.Sprintf(FileNameFormat, unitName)
}

func UnitNameByFilename(filename string) (string, error) {
	matches := FileNameRegexp.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", ErrInvalidFilename
	}
	return matches[1], nil
}
