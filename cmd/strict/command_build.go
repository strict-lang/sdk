package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/pkg/compiler"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/report"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var buildCommand = &cobra.Command{
	Use:   "build",
<<<<<<< HEAD
	Short: "Builds a Strict module",
=======
	Short: "Builds a Strict package",
>>>>>>> feature/build
	Long:  `Build compiles a file to a specified output file.`,
	RunE:   RunCompile,
}

var buildOptions struct {
	outputPath string
	reportFormat string
	debug bool
}

func init() {
	flags := buildCommand.Flags()
	flags.StringVarP(&buildOptions.outputPath, "destination", "d", "build/silk", "build destination")
	flags.BoolVarP(&buildOptions.debug, "debug", "z", false, "enable debug mode")
	flags.StringVarP(&buildOptions.reportFormat, "report-format", "r", "text",
		"format in which the report is encoded (json/pretty-json/xml/pretty-xml/text)")
}

func disableLogging() {
	if !buildOptions.debug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}

var reportFormats = map[string] func(report.Report) report.Output {
	"text": report.NewRenderingOutput,
	"json": func(input report.Report) report.Output {
		return report.NewSerializingOutput(report.NewJsonSerializationFormat(), input)
	},
	"pretty-json": func(input report.Report) report.Output {
		return report.NewSerializingOutput(report.NewPrettyJsonSerializationFormat(), input)
	},
	"xml": func(input report.Report) report.Output {
		return report.NewSerializingOutput(report.NewXmlSerializationFormat(), input)
	},
	"pretty-xml": func(input report.Report) report.Output {
		return report.NewSerializingOutput(report.NewPrettyXmlSerializationFormat(), input)
	},
}

func createFailedReport(beginTime time.Time) report.Report {
	return report.Report{
		Success:     false,
		Time:        report.Time{
			Begin: beginTime.UnixNano(),
			Completion: time.Now().UnixNano(),
		},
		Diagnostics: []report.Diagnostic{},
	}
}

func RunCompile(command *cobra.Command, arguments []string) error {
	disableLogging()
	compilationReport := compile(command, arguments)
	output := createOutput(compilationReport)
	return output.Print(command.OutOrStdout())
}

func createOutput(compilationReport report.Report) report.Output {
	if output, ok := reportFormats[buildOptions.reportFormat]; ok {
		return output(compilationReport)
	}
	return report.NewRenderingOutput(compilationReport)
}

func compile(command *cobra.Command, arguments []string) report.Report {
	beginTime := time.Now()
	file, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return createFailedReport(beginTime)
	}
	defer file.Close()
	name, err := ParseUnitName(file.Name())
	if err != nil {
		command.Printf("Invalid filename: %s\n", file.Name())
		return createFailedReport(beginTime)
	}
	return runCompilation(command, name, file)
}

func runCompilation(command *cobra.Command, unitName string, file *os.File) report.Report {
	compilation := &compiler.Compilation{
		Name:    unitName,
		Source:  &compiler.FileSource{File: file},
	}
	result := compilation.Compile()
	if result.Error != nil {
		return result.Report
	}
	if err := writeGeneratedSources(result); err != nil {
		command.PrintErrf("failed to write generated sources %v\n", err)
		result.Report.Success = false
	}
	return result.Report
}

func writeGeneratedSources(compilation compiler.Result) (err error) {
	for _, generated := range compilation.GeneratedFiles {
		if err = writeGeneratedSourceFile(generated); err != nil {
			return err
		}
	}
	return nil
}

func writeGeneratedSourceFile(generated backend.GeneratedFile) error {
	file, err := targetFile(generated.Name)
	if err != nil {
		return err
	}
	_, err = file.Write(generated.Content)
	return err
}

func targetFile(name string) (*os.File, error) {
	if buildOptions.outputPath != "" {
		return createNewFile(fmt.Sprintf("./%s", buildOptions.outputPath))
	}
	return createNewFile(fmt.Sprintf("./%s", name))
}
