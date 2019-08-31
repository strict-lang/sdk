package main

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	// ErrInvalidFieldName is returned by the ParseUnitName function when the
	// passed filename is not valid and could not be parsed.
	ErrInvalidFileName = errors.New("invalid filename")
	// TargetFileFormat is the format used to generate the name of an
	// target file. The single argument is the name of the unit.
	TargetFileFormat = "%s.cc"
	// SourceFilePattern is the regexp pattern used to parsing the unit-name
	// from a filename. It can also be used to check whether a filename is valid.
	SourceFilePattern = regexp.MustCompile(`(?P<Unit>[\w_-]+).strict`)
)

// ParseUnitName parses a unit name from the passed filename. It will
// extract the filenames actual name and ignore its file-extension. An
// error is returned when the filename does not end with the Strict
// extension or if it contains invalid characters.
func ParseUnitName(filename string) (string, error) {
	matches := SourceFilePattern.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return "", ErrInvalidFileName
	}
	return matches[1], nil
}

// GeneratedFileName returns the name of the executable that is
// generated from the unitName.
func GeneratedFileName(unitName string) string {
	return fmt.Sprintf(TargetFileFormat, unitName)
}
