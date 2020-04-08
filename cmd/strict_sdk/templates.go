package main

import "fmt"

type Template struct {
	Name string
	Directory string
	Creator func(*Make) string
}

func createStaticFile(content string) func(*Make) string {
	return func(options *Make) string {
		return content
	}
}

var readmeTemplate = Template{
	Name:      "README.md",
	Creator: func(options *Make) string {
		return fmt.Sprintf(
			`# Strict Development Kit %s`,
			options.version)
	},
}