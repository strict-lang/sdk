package main

import "flag"

var sourceFilePath string

func init() {
	flag.StringVar(&sourceFilePath, "source", "main.strict", "path to the source-file")
}

func main() {
	flag.Parse()
}
