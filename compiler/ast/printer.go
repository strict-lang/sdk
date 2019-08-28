package ast

import "github.com/tonnerre/golang-pretty"

func Print(node Node) {
	_, _ = pretty.Println(node)
}
