package main

import (
	"bufio"
	"errors"
	"github.com/urfave/cli"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/codegen"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"log"
	"os"
	"os/exec"
)

var (
	ErrNoSuchFile         = errors.New("no file with the passed name was found")
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
	unitName, err := ParseUnitName(filename)
	if err != nil {
		return err
	}
	return compileFileToDirectory(unitName, file, targetDirectory)
}

func compileFileToDirectory(unitName string, file *os.File, targetDirectory string) error {
	recorder := diagnostic.NewRecorder()
	reader := source.NewStreamReader(bufio.NewReader(file))
	log.Println("starting to parse the file")
	result := parseFile(unitName, recorder, reader)
	defer logDiagnostics(recorder, result.lines.PositionAtOffset)
	if result.err != nil {
		return reportFailedCompilation(result.err)
	}
	return generateTranslationUnit(result.parsedUnit, targetDirectory)
}

func generateTranslationUnit(unit *ast.TranslationUnit, targetDirectory string) error {
	targetFileName := codegen.FilenameByUnitName(unit.Name())
	err := generateCodeToFile(codegen.NewCodeGenerator(unit), targetFileName, targetDirectory)
	if err != nil {
		return err
	}
	return generateExecutable(targetFileName, targetDirectory)
}

func reportFailedCompilation(err error) error {
	log.Fatalf("failed to compile file: %s", err.Error())
	return ErrCompilationFailure
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

