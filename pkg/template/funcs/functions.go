package funcs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
)

const (
	FUNC_NAME_ENV            = "env"
	FUNC_NAME_STACK_OUTPUT   = "stackOutput"
	FUNC_NAME_AWS_ACCOUNT_ID = "awsAccountId"
	FUNC_NAME_HASH           = "hash"
)

// Returns empty string
func EmptyStr(s string) string {
	return ""
}

// Returns empty string without input
func EmptyInput() string {
	return ""
}

// Returns AWS account id
func AwsAccountId() (string, error) {
	c := ctlaws.NewSts(sts.New(ctlaws.AWSSess))
	output, err := c.GetCallerId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", *output.Account), nil
}

// Returns md5 hashed string
func Md5(s string) (string, error) {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", h), nil
}

// Parse environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}

// Parse cloudformation stack output
func GetStackOutputs(params ...string) (string, error) {
	if len(params) < 2 {
		return "", errors.New("Missing stack name or output key.")
	}

	name := params[0]
	key := params[1]

	c := ctlaws.NewStack(cf.New(ctlaws.AWSSess))
	if len(params) == 3 {
		c = ctlaws.NewStack(cf.New(ctlaws.GetSessionWithProfile(params[2])))
	}

	stack, err := c.DescribeStack(name)
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
