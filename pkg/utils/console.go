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

// Print format types
type FormatType string

const (
	FormatYaml FormatType = "yaml"
	FormatJson FormatType = "json"
	FormatCmd  FormatType = "cmd"
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
func StdoutStrFactory(format FormatType) StdoutStrFn {
	switch format {
	case FormatYaml:
		return toYamlStr
	case FormatJson:
		return toJsonStr
	case FormatCmd:
		return toDefault
	default:
		return toDefault
	}
}

// Print given interface to given format
func Print(format FormatType, s ...interface{}) error {
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

// Command line result print to console. It takes cmd options.
func CmdPrint(opt map[string]interface{}, format FormatType, s ...interface{}) error {
	for k, v := range opt {
		switch vv := v.(type) {
		case bool:
			switch k {
			case "quiet":
				if vv {
					return nil
				}
			}
		}
	}

	return Print(format, s...)
}
