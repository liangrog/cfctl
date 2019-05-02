package aws

import (
	"bytes"
	"fmt"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Create valid AWS s3 bucket url
func S3Url(bucket string, parts ...string) (string, error) {
	var s3url string

	base, err := url.Parse(fmt.Sprintf("https://%s.s3.amazonaws.com", bucket))
	if err != nil {
		return "", err
	}

	if len(parts) > 0 {
		if sp, err := url.Parse(path.Join(parts...)); err != nil {
			return "", err
		} else {
			s3url = base.ResolveReference(sp).String()
		}
	}

	return s3url, nil
}

// S3 wrapper
type CfS3 struct {
	Client   s3iface.S3API
	Uploader *s3manager.Uploader
}

// S3 wrapper constructor
func NewS3(s3api s3iface.S3API) *CfS3 {
	return &CfS3{Client: s3api}
}

// Check if bucket exist
func (c *CfS3) IfBucketExist(bucket string) (bool, error) {
	var err error

	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err = c.Client.HeadBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, "NotFound":
				return false, nil
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				return true, nil
			default:
				err = aerr
			}
		}

		return false, err
	}

	return true, nil
}

// Wrapper function for creatinng S3 bucket.
func (c *CfS3) CreateBucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	result, err := c.Client.CreateBucket(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return result, aerr
		}
	}

	return result, err
}

// Upload to S3 using s3manager uploader
func (c *CfS3) Upload(bucket, keyName string, body []byte, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	// Load s3manager
	if c.Uploader == nil {
		c.Uploader = s3manager.NewUploaderWithClient(c.Client, options...)
	}

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyName),
		Body:   bytes.NewReader(body),
	}

	return c.Uploader.Upload(upParams)
}
