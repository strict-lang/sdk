package diagnostic

import (
	"github.com/fatih/color"
	"io"
	"log"
)

type Kind struct {
	Name        string
	Color       *color.Color
	Description string
}

var (
	Error = Kind{
		Name:        "Error",
		Color:       color.New(color.FgRed),
		Description: "compilation can not succeed",
	}
	Info = Kind{
		Name:        "Info",
		Color:       color.New(color.FgGreen),
		Description: "compilation information",
	}
	Warning = Kind{
		Name:        "Warning",
		Color:       color.New(color.FgYellow),
		Description: "compilation can still succeed",
	}
)

func (kind Kind) Write(writer io.Writer) {
	_, err := kind.Color.Fprint(writer, kind.Name)
	if err != nil {
		log.Fatal(err)
	}
}
