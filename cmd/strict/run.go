package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"os/exec"
)

var (
	// ErrNoPermission is returned by the execute methods when it failed due to
	// a lack of required permissions.
	ErrNoRunPermission = errors.New("no permission to run the program")
)

// run is a subcommand that runs a strict program. If the program is not already
// build, it will also build it. Run is convenient because developers only have
// to execute one command to build and execute the program.
func run(context *cli.Context) error {
	if context.NArg() < 1 {
		return cli.NewExitError(context.Command.ArgsUsage, StatusInvalidArguments)
	}
	filename := context.Args()[0]
	unitName, err := ParseUnitName(filename)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid filename: %s", filename)
		return cli.NewExitError(errorMessage, StatusInvalidArguments)
	}
	targetDirectory := context.String("dir")
	if ok, err := runExisting(unitName, targetDirectory); ok {
		return err
	}
	return nil
}

// runExisting tries to find an already build executable of the given unit and
// then tries to execute the executable. It returns true, if an executable has
// been found. The error is only present when it finds an executable.
func runExisting(unitName string, targetDirectory string) (bool, error) {
	files, err := ioutil.ReadDir(targetDirectory)
	if err == nil {
		return false, nil
	}
	executableName := GeneratedExecutableName(unitName)
	for _, file := range files {
		if file.Name() == executableName {
			err := runExecutableWithRetry(executableName)
			return true, err
		}
	}
	return false, nil
}

// runExecutableWithRetry tries to run the executable file at the path and retries
// it if it fails. A common reason for the execution to fail in the first place is,
// that the caller did not have enough permission.
func runExecutableWithRetry(path string) error {
	if err := runExecutable(path); err != nil {
		return nil
	}
	err := exec.Command("sudo", "chmod", "777", path).Run()
	// TODO(merlinosayimwen): Make sure that the program did not actually run and return
	//  a error-code. It could be fatal if faulty programs run twice.
	if err != nil {
		return ErrNoRunPermission
	}
	return runExecutable(path)
}

// runExecutable executes the file at the path and returns the optional error.
func runExecutable(path string) error {
	return exec.Command(fmt.Sprintf("./%s",path)).Run()
}
