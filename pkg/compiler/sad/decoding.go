package sad

import "log"

type decoding struct {
	input  []rune
	offset int
}

func (decoding *decoding) decode() {
}

func (decoding *decoding) decodeMethod() Method {
	name := decoding.readIdentifier()
	decoding.skip('.')
	parameters := decoding.decodeParameterList()
	decoding.skip('.')
	returnType := decoding.decodeClassName()
	return Method{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
	}
}

func (decoding *decoding) decodeField() Field {
	name := decoding.readIdentifier()
	decoding.skip('.')
	className := decoding.decodeClassName()
	return Field{
		Name:  name,
		Class: className,
	}
}

func (decoding *decoding) decodeClass() Class {
	name := decoding.readIdentifier()
	methods, fields := decoding.decodeMappedClassItems()
	return Class{
		Kind:       0,
		Traits:     nil,
		Name:       name,
		Parameters: nil,
		Methods:    methods,
		Fields:     fields,
	}
}

func (decoding *decoding) decodeParameters() []Parameter {
	if decoding.isLookingAt(itemSeparator) {
		decoding.offset++
		return []Parameter{}
	}
	parameters := []Parameter{decoding.decodeParameter()}
	for decoding.isLookingAt(itemSeparator) {
		decoding.offset++
		parameters = append(parameters, decoding.decodeParameter())
	}
	return parameters
}

func (decoding *decoding) decodeTraits() []ClassName {
	if decoding.isLookingAt(itemSeparator) {
		decoding.offset++
		return []ClassName{}
	}
	traits := []ClassName{decoding.decodeClassName()}
	for decoding.isLookingAt(itemSeparator) {
		decoding.offset++
		traits = append(traits, decoding.decodeClassName())
	}
	return traits
}

func (decoding *decoding) decodeMappedClassItems() (map[string]Method, map[string]Field) {
	methods := map[string]Method{}
	fields := map[string]Field{}
	for _, item := range decoding.decodeClassItems() {
		if method, isMethod := item.(Method); isMethod {
			methods[method.Name] = method
		}
		if field, isField := item.(Field); isField {
			fields[field.Name] = field
		}
	}
	return methods, fields
}

func (decoding *decoding) decodeClassItems() []interface{} {
	if decoding.isLookingAt(classSeparator) {
		return []interface{}{}
	}
	items := []interface{}{decoding.decodeClassItem()}
	for decoding.isLookingAt(itemSeparator) {
		decoding.skip(itemSeparator)
		if decoding.isLookingAt(classSeparator) {
			break
		}
		items = append(items, decoding.decodeClassItem())
	}
	return items
}

func (decoding *decoding) decodeClassItem() interface{} {
	if decoding.isLookingAt(methodBeginKey) {
		return decoding.decodeMethod()
	}
	if decoding.isLookingAt(fieldBeginKey) {
		return decoding.decodeField()
	}
	log.Printf("invalid item specifier: '%c'", decoding.current())
	return nil
}

func (decoding *decoding) skip(value rune) {
	if !decoding.isLookingAt(value) {
		log.Printf("could not skip '%c'", value)
	}
	decoding.offset++
}

func (decoding *decoding) decodeClassName() ClassName {
	return ClassName{}
}

func (decoding *decoding) decodeParameterList() []Parameter {
	if decoding.isLookingAt(parameterListBegin) {
		return []Parameter{}
	}
	parameters := []Parameter{decoding.decodeParameter()}
	for decoding.isLookingAt(parameterSeparator) {
		parameters = append(parameters, decoding.decodeParameter())
	}
	return parameters
}

const parameterWithoutLabelIdentifierCount = 2
const parameterWithLabelIdentifierCount = 3

func (decoding *decoding) decodeParameter() Parameter {
	identifiers := decoding.readMultipleIdentifiers()
	switch len(identifiers) {
	case parameterWithLabelIdentifierCount:
	case parameterWithoutLabelIdentifierCount:
	}
	return Parameter{}
}

func (decoding *decoding) readMultipleIdentifiers() []string {
	identifiers := []string{decoding.readIdentifier()}
	for !decoding.isLookingAt('.') {
		decoding.offset++
		identifiers = append(identifiers, decoding.readIdentifier())
	}
	return identifiers
}

func (decoding *decoding) isLookingAt(value rune) bool {
	return decoding.current() == value
}

func (decoding *decoding) readIdentifier() string {
	if !isIdentifierHead(decoding.current()) {
		return ""
	}
	begin := decoding.offset
	decoding.offset++
	for isIdentifierTail(decoding.current()) {
		decoding.offset++
	}
	return decoding.extractRange(begin, decoding.offset)
}

func (decoding *decoding) extractRange(begin int, end int) string {
	span := decoding.input[begin:end]
	return string(span)
}

func (decoding *decoding) current() rune {
	return decoding.input[decoding.offset]
}

func isIdentifierHead(value rune) bool {
	return isInRange(value, 'a', 'z') ||
		isInRange(value, 'A', 'Z')
}

func isIdentifierTail(value rune) bool {
	return isIdentifierHead(value) ||
		isInRange(value, '0', '9')
}

func isInRange(character rune, begin rune, end rune) bool {
	return character >= begin && character <= end
}
