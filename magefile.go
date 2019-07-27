// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

const BinaryName = "strict"

func Build() error {
	mg.Deps(InstallDeps)
	cmd := exec.Command("go", "build", "-o", pathToBinary(), "./cmd/strict")
	return cmd.Run()
}

func Install() error {
	mg.Deps(Build)
	command := exec.Command("go", "install", ".cmd/strict")
	return command.Run()
}

func InstallDeps() error {
	cmd := exec.Command("glide")
	return cmd.Run()
}

func Clean() (err error) {
	if err = runGoClean(); err != nil {
		return
	}
	if err = deleteDirectory(pathToBinary()); err != nil {
		return
	}
	return nil
}

func deleteDirectory(path string) error {
	command := exec.Command("rm", "-rf", path)
	return command.Run()
}

func runGoClean() error {
	command := exec.Command("go", "clean")
	return command.Run()
}

func pathToBinary() string {
	goBin := os.Getenv("GOBIN")
	if goBin == "" {
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			mg.Fatal(1, "GOPATH or GOBIN not set")
			panic("no GOPATH or GOBIN set")
		}
	}
	return fmt.Sprintf("%s/%s", goBin, BinaryName)
}
