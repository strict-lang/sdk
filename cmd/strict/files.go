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

func createNewFile(filename, directory string) (*os.File, error) {
	filepath := createFilepath(filename, directory)
	if err := deleteIfExists(filepath); err != nil {
		return nil, err
	}
	if directory != "" {
		if err := createDirectoryIfNotExists(directory); err != nil {
			return nil, err
		}
	}
	return os.Create(filepath)
}

func createDirectoryIfNotExists(directory string) error {
	if _, err := os.Stat(directory); err != nil {
		return nil
	}
	dir, err := os.Create(directory)
	if err != nil {
		dir.Close()
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
