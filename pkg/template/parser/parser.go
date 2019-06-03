package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
)

const (
	FUNC_S3URL        = "tpl"
	FUNC_ENV          = "env"
	FUNC_STACK_OUTPUT = "stackOutput"
)

// Parse template by given function map and key values
func parse(s string, funcMap template.FuncMap, kv map[string]string) (bytes.Buffer, error) {
	var b bytes.Buffer

	tmpl, err := template.New(uuid.New().String()).Funcs(funcMap).Parse(s)
	if err != nil {
		return b, err
	}

	if err := tmpl.Execute(&b, kv); err != nil {
		return b, err
	}

	return b, nil
}

// Parsing template twice giving its ability to
// allow using function as value.
func doubleParse(s string, funcMap template.FuncMap, kv map[string]string) (bytes.Buffer, error) {
	b, err := parse(s, funcMap, kv)
	if err != nil {
		return b, err
	}

	//fmt.Println(funcMap)
	// Do another parse in case there are function as value
	b, err = parse(b.String(), funcMap, kv)
	if err != nil {
		return b, err
	}
	//fmt.Println(b.String())
	return b, nil
}

// Search template if it has dependency on other stacks
func SearchDependancy(s string, kv map[string]string) ([]string, error) {
	var p []string

	funcParentStack := func(name, key string) string {
		p = append(p, name)
		return ""
	}

	// Dummy
	funcS3URL := func(d string) string { return "" }
	funcEnv := func(d string) string { return "" }

	funcMap := template.FuncMap{
		FUNC_S3URL:        funcS3URL,
		FUNC_ENV:          funcEnv,
		FUNC_STACK_OUTPUT: funcParentStack,
	}

	if _, err := doubleParse(s, funcMap, kv); err != nil {
		return nil, err
	}

	return p, nil
}

// Parse template with given key-value pairs, environment variables,
// s3 template URL and stack outputs.
func Parse(s string, kv map[string]string, dc *conf.DeployConfig) ([]byte, error) {
	// Convert a give templat
	// file path to s3 url
	cfs3 := ctlaws.NewS3(s3.New(ctlaws.AWSSess))
	funcS3URL := func(path string) (string, error) {
		content, err := ioutil.ReadFile(dc.GetTplPath(path))
		if err != nil {
			return "", err
		}

		// Upload all nested template to s3
		result, err := cfs3.Upload(dc.S3Bucket, dc.GetTplPath(path), content)
		if err != nil {
			return "", err
		}

		fmt.Printf(
			"[ s3 | upload ] template: %s\tURL: %s\n",
			path,
			result.Location,
		)

		return result.Location, nil
	}

	// Parse environment variable
	funcEnv := func(key string) string {
		return os.Getenv(key)
	}

	// Parse cloudformation stack output
	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))
	funcStackOutput := func(name, key string) (string, error) {
		stack, err := stack.DescribeStack(name)
		if err != nil {
			return "", err
		}

		for _, out := range stack.Outputs {
			// Check both key and export name
			if *out.OutputKey == key || (out.ExportName != nil && *out.ExportName == key) {
				fmt.Printf(
					"[ stack | stack-output ] name: %s\tkey: %s\tvalue: %s\n",
					name,
					key,
					*out.OutputValue,
				)

				return *out.OutputValue, nil
			}
		}

		return "", errors.New(fmt.Sprintf("There is no output key %s in stack %s.", key, name))
	}

	funcMap := template.FuncMap{
		FUNC_S3URL:        funcS3URL,
		FUNC_ENV:          funcEnv,
		FUNC_STACK_OUTPUT: funcStackOutput,
	}

	output, err := doubleParse(s, funcMap, kv)

	return output.Bytes(), err
}
