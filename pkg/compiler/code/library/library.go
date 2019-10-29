package library

import "gitlab.com/strict-lang/sdk/pkg/compiler/code"

type Library interface {
	Import() (*code.Namespace,)
}

type SourceLibrary struct {}
