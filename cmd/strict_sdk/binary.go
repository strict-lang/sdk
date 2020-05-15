package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildBinary(destination string) TarballEntry {
	executionPath := selectExecutionPath()
	executeOrFail("go", "build", "-o="+destination,filepath.Join  (executionPath, "cmd/strict"))
	return TarballEntry{
		SystemPath: destination,
		Name:      filepath.Join  ("bin", selectBinaryName()),
	}
}

func selectBinaryName() string {
	if makeContext.platform == "windows" {
		return "strict.exe"
	}
	return "strict"
}

func executeOrFail(name string, options ...string) {
	command := exec.Command(name, options...)
	command.Dir = selectExecutionPath()
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		log.Fatalf("could not execute build command: %v", err.Error())
	}
}

func selectExecutionPath() string {
	if makeContext.executablePath == "." {
		return findWorkingDirectory()
	}
	if strings.HasPrefix(makeContext.executablePath, "./") {
		fixedPath := strings.Replace(makeContext.executablePath, "./", "", 1)
		return filepath.Join  (findWorkingDirectory(), fixedPath)
	}
	return makeContext.executablePath
}
