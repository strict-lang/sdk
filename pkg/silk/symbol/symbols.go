package symbol

import "fmt"

type Symbol interface {
}

type Name struct {
	Value string
}

func (name *Name) String() string {
	return fmt.Sprintf("Name{\"%s\"}", name.Value)
}
