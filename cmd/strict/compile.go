package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/codegen"
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/parser"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/source"
	"github.com/urfave/cli"
	"os"
)

var (
	ErrNoSuchFile = errors.New("no file with the passed name was found")
)

func compile(context *cli.Context) error {
	if context.NArg() < 1 {
		return cli.NewExitError(context.Command.ArgsUsage, StatusInvalidArguments)
	}
	filename := context.Args()[0]
	targetDirectory := context.String("dir")
	return compileToDirectory(filename, targetDirectory)
}

func compileToDirectory(filename string, targetDirectory string) error {
	file, err := os.Open(filename)
	if err != nil {
		return ErrNoSuchFile
	}
	defer file.Close()
	return compileFileToDirectory(filename, file, targetDirectory)
}

func compileFileToDirectory(filename string, file *os.File, targetDirectory string) error {
	unitName, err := ParseUnitName(filename)
	if err != nil {
		return err
	}
	recorder := diagnostic.NewRecorder()
	defer recorder.PrintAllEntries(diagnostic.NewFmtPrinter())

	reader := source.NewStreamReader(bufio.NewReader(file))
	tokenSource := scanner.NewDiagnosticScanner(reader, recorder)
	unit, err := parser.Parse(unitName, tokenSource, recorder)
	if err != nil {
		return err
	}
	filepath := fmt.Sprintf("%s/%s", targetDirectory, filename)
	return generateCodeToFile(codegen.NewCodeGenerator(unit), filepath)
}

func generateCodeToFile(generator *codegen.CodeGenerator, filepath string) error {
	file, err := createNewFile(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	code := generator.Generate()
	_, err = file.WriteString(code)
	return err
}

func createNewFile(filepath string) (*os.File, error) {
	if err := deleteIfExists(filepath); err != nil {
		return nil, err
	}
	return os.Create(filepath)
}

func deleteIfExists(filepath string) error {
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filepath)
}
