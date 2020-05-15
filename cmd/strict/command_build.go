package main

import (
	"github.com/spf13/cobra"
	"github.com/strict-lang/sdk/pkg/buildtool"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	"github.com/strict-lang/sdk/pkg/compiler/report"
	"io/ioutil"
	"log"
	"path/filepath"
)

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Builds a Strict package",
	Long:  `Build compiles a file to a specified output file.`,
	RunE:  RunCompile,
}

var buildOptions struct {
	outputPath   string
	reportFormat string
	backendName  string
	debug        bool
}

func init() {
	flags := buildCommand.Flags()
	flags.StringVarP(&buildOptions.backendName, "backend", "b", "c++", "backend used in code generation")
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

var reportFormats = map[string]func(report.Report, *linemap.Table) report.Output{
	"text": report.NewRenderingOutput,
	"json": func(input report.Report, lineMaps *linemap.Table) report.Output {
		return report.NewSerializingOutput(report.NewJsonSerializationFormat(), input)
	},
	"pretty-json": func(input report.Report, lineMaps *linemap.Table) report.Output {
		return report.NewSerializingOutput(report.NewPrettyJsonSerializationFormat(), input)
	},
	"xml": func(input report.Report, lineMaps *linemap.Table) report.Output {
		return report.NewSerializingOutput(report.NewXmlSerializationFormat(), input)
	},
	"pretty-xml": func(input report.Report, lineMaps *linemap.Table) report.Output {
		return report.NewSerializingOutput(report.NewPrettyXmlSerializationFormat(), input)
	},
}

func RunCompile(command *cobra.Command, arguments []string) error {
	disableLogging()
	fixOptions()
	compilationReport, lineMaps, err := runCompilation()
	if err != nil {
		return err
	}
	output := createOutput(compilationReport, lineMaps)
	return output.Print(command.OutOrStdout())
}

func fixOptions() {
	buildOptions.outputPath = filepath.Join(findWorkingDirectory(), buildOptions.outputPath)
}

func createOutput(
	compilationReport report.Report,
	table *linemap.Table) report.Output {

	if output, ok := reportFormats[buildOptions.reportFormat]; ok {
		return output(compilationReport, table)
	}
	return report.NewRenderingOutput(compilationReport, table)
}

func runCompilation() (report.Report, *linemap.Table, error) {
	directory := findWorkingDirectory()
	build := buildtool.Build{
		RootPath:      directory,
		Configuration: readBuildConfigOrFallback(directory),
	}
  return build.Run()
}

const buildFileName = `build.yml`

func readBuildConfigOrFallback(workingDirectory string) buildtool.Configuration {
	if config, err := readBuildConfig(workingDirectory); err != nil {
		return config
	}
	return buildtool.Configuration{
		PackageName:  "",
		Author:       "Undefined",
		Description:  "Undefined",
	}
}

func readBuildConfig(workingDirectory string) (buildtool.Configuration, error) {
	configPath :=filepath.Join  (workingDirectory, buildFileName)
	return buildtool.ReadConfiguration(configPath)
}
