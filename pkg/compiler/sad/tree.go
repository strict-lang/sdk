package sad

type TypeKind int8

const (
	ClassKind TypeKind = iota
	TraitKind
)

type Node interface {
	encode(encoding *encoding)
}

type Class struct {
	Kind       TypeKind
	Traits     []ClassName
	Name       string
	Parameters []TypeParameter
	Methods    map[string]Method
	Fields     map[string]Field
}

type ClassName struct {
	Name      string
	Wildcard  bool
	Arguments []ClassName
}

type ClassArgument struct {
	Class    ClassName
	Wildcard bool
}

type TypeParameter struct {
	Class    ClassName
	Wildcard bool
}

func (class *Class) FindMethod(name string) (Method, bool) {
	if method, ok := class.Methods[name]; ok {
		return method, true
	}
	//for _, super := range class.Traits {
	//  if superMethod, ok := super.FindMethod(name); ok {
	//    return superMethod, true
	//  }
	//}
	return Method{}, false
}

func (class *Class) FindField(name string) (Field, bool) {
	if field, ok := class.Fields[name]; ok {
		return field, true
	}
	return Field{}, false
}

type Method struct {
	Name       string
	Parameters []Parameter
	ReturnType ClassName
}

type Parameter struct {
	Name  string
	Label string
	Class ClassName
}

type Field struct {
	Name  string
	Class ClassName
}
