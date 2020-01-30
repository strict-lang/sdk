package ilasm

type Class struct {}

var Float = &Class{}
var Int = &Class{}
var String = &Class{}

func (class *Class) IsAssignable(target *Class) bool {
	return false
}

