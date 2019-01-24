package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type MessageType string

// Message types
const (
	MessageTypeInfo MessageType = "INFO"

	MessageTypeWarning MessageType = "WARNING"

	MessageTypeError MessageType = "ERROR"

	MessageTypeFatal MessageType = "FATAL"
)

// Function type to output to stdout
type StdoutStrFn func(input interface{}) (string, error)

// Convert to yaml string
func toYamlStr(input interface{}) (string, error) {
	output, err := yaml.Marshal(input)
	return string(output), err
}

// Convert to json string with indent
func toJsonStr(input interface{}) (string, error) {
	output, err := json.MarshalIndent(input, "", "    ")
	return string(output), err
}

// default
func toDefault(input interface{}) (string, error) {
	return fmt.Sprintf("%s", input), nil
}

// stdout string conversion factory
// Default to json
func StdoutStrFactory(format string) StdoutStrFn {
	switch format {
	case "yaml":
		return toYamlStr
	case "json":
		return toJsonStr
	default:
		return toDefault
	}
}

// Print given interface to given format
func Print(format string, s ...interface{}) error {
	if len(s) == 0 {
		return errors.New(MsgFormat("Printing output error: No object given", MessageTypeError))
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
func MsgFormat(msg string, msgType MessageType, options ...string) string {
	return fmt.Sprintf("%s", msg)
}
