package data

import (
	"encoding/json"
	"io"
)

// deserialize the object from Json string in an io.reader to a given interface
func FromJSON(in interface{}, reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(in)
}

// serialize the object to string based Json format
func ToJSON(in interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)

	return encoder.Encode(in)
}
