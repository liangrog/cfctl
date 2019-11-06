package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/template/funcs"
)

const (
	FUNC_S3URL = "tpl"
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

	funcMap := template.FuncMap{
		funcs.FUNC_NAME_STACK_OUTPUT:   funcParentStack,
		FUNC_S3URL:                     funcs.EmptyStr,
		funcs.FUNC_NAME_ENV:            funcs.EmptyStr,
		funcs.FUNC_NAME_AWS_ACCOUNT_ID: funcs.EmptyInput,
		funcs.FUNC_NAME_HASH:           funcs.EmptyStr,
	}

	if _, err := parse(s, funcMap, kv); err != nil {
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

	funcMap := template.FuncMap{
		FUNC_S3URL:                     funcS3URL,
		funcs.FUNC_NAME_ENV:            funcs.GetEnv,
		funcs.FUNC_NAME_STACK_OUTPUT:   funcs.GetStackOutputs,
		funcs.FUNC_NAME_AWS_ACCOUNT_ID: funcs.AwsAccountId,
		funcs.FUNC_NAME_HASH:           funcs.Md5,
	}

	output, err := parse(s, funcMap, kv)

	return output.Bytes(), err
}
