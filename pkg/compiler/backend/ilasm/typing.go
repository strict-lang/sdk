package ilasm

type Class struct {
	Name string
}

var Float = &Class{
	Name: "Float",
}
var Int = &Class{
	Name: "Int",
}
var String = &Class{
	Name: "String",
}

func (class *Class) IsAssignable(target *Class) bool {
	return class.Name == target.Name
}

