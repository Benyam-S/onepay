package tools

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
)

// Marshal is a function that marshal a given data to a needed format
func Marshal(data interface{}, format string) ([]byte, error) {

	if format == "xml" {

		dataType := reflect.TypeOf(data)

		// since xml doesn't marshal maps, we made our's
		if dataType.String() == "map[string]string" {

			obj, ok := data.(map[string]string)
			if !ok {
				return nil, errors.New("unable to parse map")
			}

			output := "<result>"
			for key, value := range obj {
				formatableString := "<%s>%s</%s>"
				output = fmt.Sprintf(output+formatableString, key, value, key)
			}
			output += "</result>"

			fmt.Println(output)
			return []byte(output), nil
		}

		output, err := xml.Marshal(data)
		if err != nil {
			panic(err)
		}

		return output, nil
	}

	// json is the standard marshal unit
	return json.Marshal(data)
}

// MarshalIndent is a function that marshal a given data with indentation to a needed format
func MarshalIndent(data interface{}, prefix, identation string, format string) ([]byte, error) {

	if format == "xml" {

		dataType := reflect.TypeOf(data)

		// since xml doesn't marshal maps, we made our's
		if dataType.String() == "map[string]string" {

			obj, ok := data.(map[string]string)
			if !ok {
				return nil, errors.New("unable to parse map")
			}

			output := prefix + "<result>\n"
			for key, value := range obj {
				formatableString := "%s%s<%s>%s</%s>\n"
				output = fmt.Sprintf(output+formatableString, prefix, identation, key, value, key)
			}
			output += prefix + "</result>"

			return []byte(output), nil
		}

		output, err := xml.MarshalIndent(data, prefix, identation)
		if err != nil {
			panic(err)
		}

		return output, nil
	}

	// json is the standard marshal unit
	return json.MarshalIndent(data, prefix, identation)
}
