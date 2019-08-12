package codegen

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	FileNameFormat     = "%s.generated.cc"
	ArduinoFileNameFormat = "%s.ino"
	FileNameRegexp     = regexp.MustCompile(`(?P<Unit>\w+)\.generated\.(ino|cc)`)
	ErrInvalidFilename = errors.New("invalid filename format")
)

// FilenameByUnitName returns the name of a file that belongs to the passed
// translation unit.
func (generator *CodeGenerator) Filename() string {
	if generator.settings.IsTargetingArduino {
		return fmt.Sprintf(ArduinoFileNameFormat, generator.unit.Name())
	}
	return fmt.Sprintf(FileNameFormat, generator.unit.Name())
}

func UnitNameByFilename(filename string) (string, error) {
	matches := FileNameRegexp.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", ErrInvalidFilename
	}
	return matches[1], nil
}
