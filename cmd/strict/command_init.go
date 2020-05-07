package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize strict package",
	Long:  `init generates a set of files that should be in every strict package`,
	RunE:  RunInit,
}

var initOptions struct {
	name        string
	author      string
	description string
}

func init() {
	flags := initCommand.Flags()
	flags.StringVarP(&initOptions.name, "name", "n", "", "package name")
	flags.StringVarP(&initOptions.author, "author", "a", "", "package author")
	flags.StringVarP(&initOptions.description, "description", "d", "None", "package description")
}

func RunInit(command *cobra.Command, arguments []string) error {
	workingDirectory := findWorkingDirectory()
	if initOptions.name == "" {
		fetchOptionsInteractive(workingDirectory)
	}
	if err := createPackage(workingDirectory); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to init package %s: %s", initOptions.name, err.Error())
		return err
	}
	fmt.Println("Successfully created package")
	return nil
}

const noDescription = "None"

func fetchOptionsInteractive(workingDirectory string) {
	inputReader := bufio.NewReader(os.Stdin)
	fallbackPackageName := strings.Title(strings.ToLower(filepath.Base(workingDirectory)))
	initOptions.name = promptWithFallback(inputReader, "Package name", fallbackPackageName)
	initOptions.author = promptWithFallback(inputReader, "Package author", findFallbackName())
	initOptions.description = promptWithFallback(inputReader, "Package description", noDescription)
}

func findFallbackName() string {
	command := exec.Command("sh", "-c", "git config --global --get user.name")
	rawMessage, err := command.CombinedOutput()
	if err != nil {
		return findOperatingSystemName()
	}
	message := string(rawMessage)
	if message == "" {
		return findOperatingSystemName()
	}
	return strings.ReplaceAll(message, "\n", "")
}

const fallbackName = "Unknown"

func findOperatingSystemName() string {
	current, err := user.Current()
	if err != nil {
		return fallbackName
	}
	if current.Name == "" {
		return resolveBaseUserName(current.Username)
	}
	return current.Name
}

func resolveBaseUserName(name string) string {
	slash := strings.Index(name, "\\")
	if slash >= 0 {
		rightSideStart := slash + 1
		if rightSideStart < len(name) {
			return name[rightSideStart:]
		}
	}
	return name
}

func promptWithFallback(input *bufio.Reader, text string, fallback string) string {
	fmt.Printf("%s (%s): ", text, fallback)
	rawAnswer, _, err := input.ReadLine()
	if err != nil {
		return fallback
	}
	answer := string(rawAnswer)
	if answer == "" {
		return fallback
	}
	return answer
}

func createPackage(root string) error {
	initializeGit()
	for _, file := range templateFiles {
		if err := file.Save(root); err != nil {
			return err
		}
	}
	sourcePath := root + "/src"
	_, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return os.Mkdir(sourcePath, simplePermissions)
	}
	return err
}

func initializeGit() {
	command := exec.Command("sh", "-c", "git init")
	command.Stderr = ioutil.Discard
	command.Stdout = ioutil.Discard
	_ = command.Run()
}

type Template struct {
	content func() string
	name    string
}

var templateFiles = []Template{
	gitignore,
	editorconfig,
	buildConfig,
}

const simplePermissions = 0644

func (template *Template) Save(root string) error {
	encodedContent := []byte(template.content())
	path := root + "/" + template.name
	return ioutil.WriteFile(path, encodedContent, simplePermissions)
}

func staticContent(input string) func() string {
	return func() string {
		return input
	}
}

const editorconfigContent = `root = true

[*]
indent_size = 2
indent_style = tab
tab_width = 2
trim_trailing_whitespace = true
insert_final_newline = false
charset = utf-8
max_line_length = 80
end_of_line = lf
`

var editorconfig = Template{
	name:    ".editorconfig",
	content: staticContent(editorconfigContent),
}

const gitignoreContent = `*.exe
*.dll
*.silk

.build/
`

var gitignore = Template{
	name:    ".gitignore",
	content: staticContent(gitignoreContent),
}

var buildConfig = Template{
	name:    "build.yml",
	content: createBuildConfig,
}

const baseBuildConfig = `name: %s
author: %s`

func createBuildConfig() string {
	base := fmt.Sprintf(baseBuildConfig, initOptions.name, initOptions.author)
	if hasDescription() {
		return base + "\ndescription: " + initOptions.description
	}
	return base
}

func hasDescription() bool {
	description := initOptions.description
	return description != "" && description != noDescription
}
