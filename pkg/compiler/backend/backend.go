package backend

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/isolate"
	"io/ioutil"
	"os"
	"path/filepath"
)

type GeneratedFile struct {
	Name    string
	Content []byte
}

func (file *GeneratedFile) Save(rootDirectory string) error {
	fullPath := filepath.Join(rootDirectory, file.Name)
	parent := filepath.Dir(fullPath)
	if err := createDirectoryIfNotExists(parent); err != nil {
		return fmt.Errorf("could not create directory: %v", err)
	}
	return ioutil.WriteFile(fullPath, file.Content, os.ModePerm)
}

func createDirectoryIfNotExists(directory string) error {
	err := os.MkdirAll(directory, os.ModePerm)
	if os.IsExist(err) {
		return nil
	}
	return err
}

type Output struct {
	GeneratedFiles []GeneratedFile
}

type Input struct {
	Unit        *tree.TranslationUnit
	Diagnostics *diagnostic.Bag
}

type Backend interface {
	Generate(Input) (Output, error)
}

func LookupInIsolate(isolate *isolate.Isolate, name string) (Backend, bool) {
	propertyKey := createPropertyKey(name)
	if property, ok := isolate.Properties.Lookup(propertyKey); ok {
		if factory, ok := property.(func() Backend); ok {
			return factory(), true
		}
	}
	return nil, false
}

func createPropertyKey(name string) string {
	return "backend." + name
}

func Register(name string, factory func() Backend) {
	isolate.RegisterConfigurator(func(isolate *isolate.Isolate) {
		propertyKey := createPropertyKey(name)
		isolate.Properties.Insert(propertyKey, factory)
	})
}
