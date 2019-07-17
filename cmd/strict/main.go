package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	StatusInvalidArguments = 80
	StatusFileNotFound     = 81
	StatusInvalidFile      = 82
	StatusBuildFailed      = 83
	StatusNoPermission     = 84
)

func main() {
	app := cli.NewApp()
	app.Name = "strict"
	app.HelpName = "strict"
	app.Usage = "strict's build tool"
	app.Description = `The strict cli-tool helps developers to build and execute strict programs`
	app.Version = "0.2.1"

	app.Commands = []cli.Command{
		{
			Name:      "compile",
			Aliases:   []string{"c"},
			Usage:     "compile a strict file",
			Action:    compile,
			ArgsUsage: "compile <path>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "the target directory",
					Value: "",
				},
			},
		},
		{
			Name: "format",
			Aliases: []string{"f"},
			Usage: "format a strict file",
			Action: formatCode,
			ArgsUsage: "format <path> -o -t=target.strict",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "target, t",
					Usage: "The target file",
					Value: "",
				},
				cli.BoolFlag{
					Name: "override, o",
					Usage: "Sets whether file is overridden",
				},
			},
		},
		{
			Name:      "run",
			Aliases:   []string{"r"},
			Usage:     "run a strict program",
			Action:    run,
			ArgsUsage: "run <path>",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
