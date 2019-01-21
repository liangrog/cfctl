package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Function type to output to stdout
type StdoutStrFn func(input interface{}) (string, error)

// Convert to yaml string
func ToYamlStr(input interface{}) (string, error) {
	output, err := yaml.Marshal(input)
	return string(output), err
}

// Convert to json string with indent
func ToJsonStr(input interface{}) (string, error) {
	output, err := json.MarshalIndent(input, "", "    ")
	return string(output), err
}

// stdout string conversion factory
// Default to json
func StdoutStrFactory(format string) StdoutStrFn {
	switch format {
	case "yaml":
		return ToYamlStr
	default:
		return ToJsonStr
	}
}

// Print given interface to given format
func Print(format string, s ...interface{}) error {
	if len(s) == 0 {
		return errors.New(MsgFormat("Printing output error: No object given"))
	}

	fn := StdoutStrFactory(format)
	for _, ss := range s {
		out, err := fn(ss)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	return nil
}

// Format error message
func MsgFormat(msg string, options ...string) string {
	return fmt.Sprintf("%s", msg)
}
