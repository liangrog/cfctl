package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
)

// mock client
type stsFakeClient struct {
	stsiface.STSAPI
}

var fsts = NewSts(&stsFakeClient{})

func (fc *stsFakeClient) GetCallerIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	result := &sts.GetCallerIdentityOutput{
		Account: aws.String("123456"),
		Arn:     aws.String("arn:aws:iam::123456:user/Bob"),
		UserId:  aws.String("12345"),
	}
	return result, nil
}

func TestGetCallerId(t *testing.T) {
	result, err := fsts.GetCallerId()
	assert.Equal(t, "123456", *result.Account)
	assert.NoError(t, err)
}
