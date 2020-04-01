package report

import "io"

type Output interface {
	Print(writer io.Writer) error
}

type serializingOutput struct {
	format SerializationFormat
	report Report
}

func NewSerializingOutput(format SerializationFormat, report Report) Output {
	return serializingOutput{
		format: format,
		report: report,
	}
}

func (output serializingOutput) Print(writer io.Writer) error {
	encoded, err := output.format.Marshal(output.report)
	if err != nil {
		return err
	}
	_, err = writer.Write(encoded)
	return err
}
