// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"runtime"
)

const BinaryName = "strict"

func binaryName() string {
	if runtime.GOOS == "windows" {
		return BinaryName + ".exe"
	}
	return BinaryName
}

func Install() error {
	mg.Deps(InstallDeps)
	if err := sh.RunV(mg.GoCmd(), "install", "./cmd/strict"); err != nil {
		fmt.Println("build - failed to install binary")
		return err
	}
	fmt.Println("build - successfully installed binaries")
	return nil
}

func InstallDeps() error {
	if err := sh.RunV("glide", "install"); err != nil {
		fmt.Println("build - failed tp install dependencies")
		return err
	}
	fmt.Println("build - successfully installed dependencies")
	return nil
}
