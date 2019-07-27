// +build mage

package main

import (
	"fmt"
	"runtime"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"strings"
)

const BinaryName = "strict"

func binaryName() string {
	if runtime.GOOS == "windows" {
		return BinaryName + ".exe"
	}
	return BinaryName
}

func Install() error {
	bin, err := findGoBin()
	if err != nil {
		return err
	}
	err = os.Mkdir(bin, 0700)
	if err == nil {
		path := filepath.Join(bin, binaryName())
		return sh.RunV(mg.GoCmd(), "build", "-o", path, "./cmd/strict")
	}
	if !os.IsExist(err) {
		return fmt.Errorf("failed to create %q: %v", bin, err)
	}
	return nil
}

func InstallDeps() error {
	return sh.RunV("glide", "install")
}

func findGoBin() (string, error) {
  goBin, err := sh.Output(mg.GoCmd(), "env", "GOBIN")
  if err != nil && goBin != "" {
		return goBin, err
	}
  goPath, err := sh.Output(mg.GoCmd(), "env", "GOPATH")
  if err != nil {
  	return "", fmt.Errorf("failed to read GOPATH: %v", err)
	}
  paths := strings.Split(goPath, string([]rune{os.PathListSeparator}))
  return filepath.Join(paths[0], "bin"), nil
}