package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var makeCommand = &cobra.Command{
	Use:   "make",
	Short: "",
	RunE:  RunMake,
}

type Make struct {
	bundle           bool
	resultOutputPath string
	workOutputPath   string
	version          string
	name             string
	platform         string
	architecture     string
	executablePath   string
}

var makeContext Make

func init() {
	flags := makeCommand.Flags()
	flags.BoolVarP(&makeContext.bundle, "bundle", "b", false, "puts the created sdk into a gzipped tarball")
	flags.StringVarP(&makeContext.architecture, "arch", "a", "x86_64", "architecture of the target platform")
	flags.StringVarP(&makeContext.platform, "platform", "p", "linux", "target platform")
	flags.StringVarP(&makeContext.resultOutputPath, "output", "o", "", "output path of the created strict_sdk")
	flags.StringVarP(&makeContext.version, "version", "v", "undefined", "version of the created strict_sdk")
	flags.StringVarP(&makeContext.executablePath, "executable", "e", ".", "path to the sdk executable")
}

func prepareOptions() {
	makeContext.name = fmt.Sprintf("sdk-%s-%s-%s.tar.gz",
		strings.ReplaceAll(makeContext.version, "/", "-"),
		makeContext.platform,
		makeContext.architecture)
	
	makeContext.resultOutputPath = path.Join(findWorkingDirectory(), fixPath(makeContext.resultOutputPath))
	makeContext.workOutputPath = chooseWorkOutputPath()
	makeContext.ensureDirectoryExists(makeContext.workOutputPath)
	makeContext.ensureDirectoryExists(makeContext.resultOutputPath)
}

func chooseWorkOutputPath() string {
	if makeContext.bundle {
		return makeContext.resultOutputPath + "/temp"
	}
	return makeContext.resultOutputPath
}

func RunMake(command *cobra.Command, arguments []string) error {
	prepareOptions()
	files := append(createTemplateFiles(), downloadBasePackage()...)
	files = append(files, TarballEntry{
		SystemPath: makeContext.executablePath,
		Name:       path.Join("bin", selectBinaryName()),
	})
	if makeContext.bundle {
		bundleSdk(files)
	}
	return nil
}

const basePackageUrl = "https://github.com/strict-lang/Strict.Base/tarball/master"
const basePackageDirectory = "library"

func downloadBasePackage() []TarballEntry {
	log.Printf("downloading base package")
	destination := makeContext.createPath(basePackageDirectory)
	makeContext.ensureDirectoryExists(destination)
	download := NewLibraryDownload(basePackageUrl, destination)
	files, err := download.Download()
	if err != nil {
		log.Fatalf("could not download base package: %v", err)
	}
	return changeParentDirectory(files, basePackageDirectory)
}

func changeParentDirectory(entries []TarballEntry, directory string) (fixed []TarballEntry) {
	for _, entry := range entries {
		fixed = append(fixed, TarballEntry{
			SystemPath: entry.SystemPath,
			Name:       path.Join(directory, entry.Name),
		})
	}
	return
}

func bundleSdk(files []TarballEntry) {
	tarball := &Tarball{
		Path:    makeContext.createResultPath(makeContext.name),
		Entries: files,
	}
	if err := tarball.Write(); err != nil {
		log.Fatalf("failed to create tarball: %v", err.Error())
	}
	_ = os.RemoveAll(makeContext.workOutputPath)
}

const basePermission = 0777

const readmeContent = `# Strict Development Kit`

const installContent = `# Installation`

const indexContent = `<!DOCTYPE html>
<html>
  <head></head>
	<body>
		<h1>Strict Development Kit</h1>
		<p>This is not really a documentation</p>
	</body>
</html>`

func createTemplateFiles() []TarballEntry {
	makeContext.ensureOutputPathExists()
	return []TarballEntry{
		makeContext.writeFile("README.md", readmeContent),
		makeContext.writeFile("INSTALL.md", installContent),
		makeContext.writeFile(makeContext.inNewDirectory("documentation", "index.html"), indexContent),
	}
}

func (context *Make) ensureOutputPathExists() {
	err := os.MkdirAll(context.workOutputPath, basePermission)
	context.ensureSuccess(context.workOutputPath, err)
}

func (context *Make) inNewDirectory(directory string, file string) string {
	context.ensureDirectoryExists(context.createPath(directory))
	return path.Join(directory, file)
}

func (context *Make) writeFile(name string, content string) TarballEntry {
	systemPath := context.createPath(name)
	err := ioutil.WriteFile(systemPath, []byte(content), basePermission)
	context.ensureSuccess(name, err)
	return TarballEntry{
		Name:       name,
		SystemPath: systemPath,
	}
}

func (context *Make) ensureDirectoryExists(directory string) {
	err := os.MkdirAll(directory, basePermission)
	context.ensureSuccess(directory, err)
}

func (context *Make) ensureSuccess(file string, err error) {
	if err != nil {
		log.Fatalf("could not write file %s: %v", file, err.Error())
	}
}

func (context *Make) createPath(name string) string {
	return path.Clean(path.Join(context.workOutputPath, name))
}

func (context *Make) createResultPath(name string) string {
	return path.Clean(path.Join(context.resultOutputPath, name))
}

func findWorkingDirectory() string {
	directory, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	return fixPath(directory)
}

func fixPath(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), "/")
}
