package compiler

import (
	"os"
)

// CompileFile compiles the passed file.
func CompileFile(unitName string, file *os.File) Result {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name:   unitName,
	}
	return compilation.Compile()
}

// CompileString compiles the passed string.
func CompileString(name string, value string) Result {
	compilation := &Compilation{
		Source: &InMemorySource{Source: value},
		Name:   name,
	}
	return compilation.Compile()
}

// ParseFile parses the passed file.
func ParseFile(unitName string, file *os.File) ParseResult {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name:   unitName,
	}
	return compilation.parse()
}
