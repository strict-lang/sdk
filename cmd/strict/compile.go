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
	"log"
	"os"
	"os/exec"
)

var (
	ErrNoSuchFile = errors.New("no file with the passed name was found")
	ErrCompilationFailure = errors.New("compilation failure")
)

func compile(context *cli.Context) error {
	if context.NArg() < 1 {
		return cli.NewExitError(context.Command.ArgsUsage, StatusInvalidArguments)
	}
	filename := context.Args()[0]
	log.Printf("starting to compile %s", filename)
	targetDirectory := context.String("dir")
	err := compileToDirectory(filename, targetDirectory)
	if err != nil {
		return err
	}
	log.Printf("successfully build %s\n", filename)
	return nil
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

	log.Println("starting to parse the file")
	reader := source.NewStreamReader(bufio.NewReader(file))
	tokenSource := scanner.NewDiagnosticScanner(reader, recorder)
	unit, err := parser.Parse(unitName, tokenSource, recorder)

	if err != nil {
		log.Fatalf("failed to compile file: %s", err.Error())
		return ErrCompilationFailure
	}
	targetFileName := codegen.FilenameByUnitName(unitName)
	err = generateCodeToFile(codegen.NewCodeGenerator(unit), targetFileName, targetDirectory)
	if err != nil {
		return err
	}
	return generateExecutable(targetFileName, targetDirectory)
}

func generateExecutable(filename, directory string) error {
	filepath := createFilepath(filename, directory)
	return exec.Command("g++", filepath).Run()
}

func generateCodeToFile(generator *codegen.CodeGenerator, filename, directory string) error {
	file, err := createNewFile(filename, directory)
	if err != nil {
		return err
	}
	defer file.Close()
	code := generator.Generate()
	writer := bufio.NewWriter(file)
	if _, err := writer.Write([]byte(code)); err != nil {
		return err
	}
	return writer.Flush()
}

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
