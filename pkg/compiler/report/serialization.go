package report

import (
	"encoding/json"
	"encoding/xml"
)

type SerializationFormat struct {
	marshal   func(value interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

const prettyJsonIndent = "  "

var prettyJsonSerializationFormat = SerializationFormat{
	marshal: func(value interface{}) ([]byte, error) {
		return json.MarshalIndent(value, "", prettyJsonIndent)
	},
	unmarshal: json.Unmarshal,
}

func NewPrettyJsonSerializationFormat() SerializationFormat {
	return prettyJsonSerializationFormat
}

var jsonSerializationFormat = SerializationFormat{
	marshal:   json.Marshal,
	unmarshal: json.Unmarshal,
}

func NewJsonSerializationFormat() SerializationFormat {
	return jsonSerializationFormat
}

const prettyXmlIndent = "  "

var prettyXmlSerializationFormat = SerializationFormat{
	marshal: func(value interface{}) ([]byte, error) {
		return xml.MarshalIndent(value, "", prettyXmlIndent)
	},
	unmarshal: xml.Unmarshal,
}

func NewPrettyXmlSerializationFormat() SerializationFormat {
	return prettyXmlSerializationFormat
}

var xmlSerializationFormat = SerializationFormat{
	marshal:   xml.Marshal,
	unmarshal: xml.Unmarshal,
}

func NewXmlSerializationFormat() SerializationFormat {
	return xmlSerializationFormat
}

func (format SerializationFormat) Marshal(report Report) ([]byte, error) {
	return format.marshal(report)
}

func (format SerializationFormat) Unmarshal(bytes []byte) (Report, error) {
	result := Report{}
	if err := format.unmarshal(bytes, &result); err != nil {
		return Report{}, err
	}
	return result, nil
}
