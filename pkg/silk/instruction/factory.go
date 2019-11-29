package instruction

type LabelValuePair struct {
	Label Label
	Value FieldOrValue
}

type FieldOrValue struct{}

func Phi(options ...LabelValuePair) Instruction {
	return Instruction{}
}
