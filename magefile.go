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
		fmt.Println("[MAGE] Failed to install binary")
		return err
	}
	fmt.Println("[MAGE] Successfully installed binaries")
	return nil
}

func InstallDeps() error {
	if err := sh.RunV("glide", "install"); err != nil {
		fmt.Println("[MAGE] Failed tp install dependencies")
		return err
	}
	fmt.Println("[MAGE] Successfully installed dependencies")
	return nil
}