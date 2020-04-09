package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

type TarballEntry struct {
	SystemPath string
	Name       string
}

type Tarball struct {
	Path    string
	Entries []TarballEntry
}

func (tarball *Tarball) Write() error {
	file, err := os.Create(tarball.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	zipWriter := gzip.NewWriter(file)
	defer zipWriter.Close()
	tarWriter := tar.NewWriter(zipWriter)
	defer tarWriter.Close()
	return tarball.writeTo(tarWriter)
}

func (tarball *Tarball) writeTo(writer *tar.Writer) error {
	for _, entry := range tarball.Entries {
		if err := tarball.writeEntry(writer, entry); err != nil {
			return err
		}
	}
	return nil
}

func (tarball *Tarball) writeEntry(writer *tar.Writer, entry TarballEntry) error {
	file, err := os.Open(entry.SystemPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return tarball.writeFile(writer, file, entry)
}

func (tarball *Tarball) writeFile(
	writer *tar.Writer,
	file *os.File,
	entry TarballEntry) error {

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := createTarHeader(entry.Name, stat)
	if err = writer.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func createTarHeader(name string, stat os.FileInfo) *tar.Header {
	return &tar.Header{
		Name:    name,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
}
