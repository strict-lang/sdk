package main

import (
	"fmt"
	"os"
)

func createFilepath(filename, directory string) string {
	if directory == "" {
		return filename
	}
	return fmt.Sprintf("%s/%s", directory, filename)
}

func createNewFile(filepath string) (*os.File, error) {
	if err := deleteIfExists(filepath); err != nil {
		return nil, err
	}
	return os.Create(filepath)
}

func createDirectoryIfNotExists(directory string) error {
	err := os.MkdirAll(directory, os.ModePerm)
	if os.IsExist(err) {
		return nil
	}
	return err
}

func deleteIfExists(filepath string) error {
	err := os.Remove(filepath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
