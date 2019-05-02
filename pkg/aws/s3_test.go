package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

// mock client
type s3FakeClient struct {
	s3iface.S3API
}

var cfs3 = NewS3(&s3FakeClient{})

func (fc *s3FakeClient) HeadBucket(input *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	if *input.Bucket == "notexist" {
		return &s3.HeadBucketOutput{}, awserr.New(s3.ErrCodeNoSuchBucket, "notexist", errors.New("test"))
	}

	return &s3.HeadBucketOutput{}, nil
}

func (fc *s3FakeClient) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return &s3.CreateBucketOutput{}, nil
}

func TestIfBucketExist(t *testing.T) {
	result, err := cfs3.IfBucketExist("notexist")
	assert.False(t, result)
	assert.NoError(t, err)

	result, err = cfs3.IfBucketExist("exist")
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestCreateBucket(t *testing.T) {
	_, err := cfs3.CreateBucket(&s3.CreateBucketInput{})
	assert.NoError(t, err)
}
