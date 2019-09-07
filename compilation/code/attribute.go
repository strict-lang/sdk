package code

type AttributeKind int

const (
	TypeAttribute AttributeKind = iota
	InstanceFieldAttribute
	LocalFieldAttribute
	MethodAttribute
)

type Attribute struct {
	Identifier string
	Kind       AttributeKind
	link       interface{}
}

func NewTypeAttribute(identifier string, linkedType *Type) *Attribute {
	return &Attribute{
		Identifier: identifier,
		Kind:       TypeAttribute,
		link:       linkedType,
	}
}

func (attribute *Attribute) LinkedType() (linkedType *Type, exists bool) {
	linkedType, exists = attribute.link.(*Type)
	return
}

func (attribute *Attribute) LinkedMethod() (linkedMethod *TypedMethod, exists bool) {
	linkedMethod, exists = attribute.link.(*TypedMethod)
	return
}

func (attribute *Attribute) LinkedNamespace() (linkedNamespace *Namespace, exists bool) {
	linkedNamespace, exists = attribute.link.(*Namespace)
	return linkedNamespace, exists
}

func (attribute *Attribute) LinkedInstanceField() (field *TypedInstanceField,
	exists bool) {
	field, exists = attribute.link.(*TypedInstanceField)
	return
}

func (attribute *Attribute) LinkedLocalField() (field *TypedLocalField, exists bool) {
	field, exists = attribute.link.(*TypedLocalField)
	return
}
