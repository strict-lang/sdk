package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LibraryDownload struct {
	destination         string
	files               []TarballEntry
	reader              *tar.Reader
	basePackageUrl      string
	ignoreFirstLayer    bool
	firstLayerDirectory string
}

func NewLibraryDownload(url string, destination string) *LibraryDownload {
	return &LibraryDownload{
		destination:      destination,
		basePackageUrl:   url,
		ignoreFirstLayer: true,
	}
}

func (download *LibraryDownload) Download() ([]TarballEntry, error) {
	response, err := http.Get(download.basePackageUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	err = download.unpack(response.Body)
	return download.files, err
}

func (download *LibraryDownload) unpack(reader io.Reader) error {
	zipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	download.reader = tar.NewReader(zipReader)
	for {
		shouldContinue, err := download.processNextEntry()
		if err != nil {
			return err
		}
		if !shouldContinue {
			break
		}
	}
	return nil
}

func (download *LibraryDownload) processNextEntry() (shouldContinue bool, err error) {
	switch header, err := download.reader.Next(); {
	case err == io.EOF:
		return false, nil
	case err != nil:
		return false, err
	case header == nil:
		return true, nil
	default:
		err = download.unpackEntry(header)
		return err == nil, err
	}
}

func (download *LibraryDownload) unpackEntry(header *tar.Header) error {
	target := filepath.Join(download.destination, download.fixPath(header.Name))
	switch header.Typeflag {
	case tar.TypeDir:
		return download.unpackDirectory(target, header)
	case tar.TypeReg:
		return download.unpackFile(target, header)
	}
	return nil
}

func (download *LibraryDownload) unpackDirectory(target string, header *tar.Header) error {
	if download.ignoreFirstLayer && download.firstLayerDirectory == "" {
		download.firstLayerDirectory = path.Clean(header.Name)
		return nil
	}
	return os.MkdirAll(target, basePermission)
}

const fileFlags = os.O_CREATE | os.O_RDWR

func (download *LibraryDownload) unpackFile(target string, header *tar.Header) error {
	if download.shouldIgnoreCurrentLayer() {
		return nil
	}
	download.addFile(target, header)
	file, err := os.OpenFile(target, fileFlags, os.FileMode(header.Mode))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, download.reader)
	return err
}

func (download *LibraryDownload) fixPath(path string) string {
	if !download.ignoreFirstLayer || download.firstLayerDirectory == "" {
		return path
	}
	return strings.Replace(path, download.firstLayerDirectory+"/", "", 1)
}

func (download *LibraryDownload) shouldIgnoreCurrentLayer() bool {
	return download.ignoreFirstLayer && download.firstLayerDirectory == ""
}

func (download *LibraryDownload) addFile(systemPath string, header *tar.Header) {
	name := download.fixPath(header.Name)
	download.files = append(download.files, TarballEntry{
		SystemPath: systemPath,
		Name:       name,
	})
}
