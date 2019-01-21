package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// Fetch AWS immediate error code and message
func AWSErrDetail(err error) (string, string, awserr.Error) {
	var code, msg string
	var awsErr awserr.Error

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			code = awsErr.Code()
			msg = awsErr.Message()
		}
	}
	return code, msg, awsErr
}
