package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/stretchr/testify/assert"
)

// mock client
type fakeClient struct {
	cloudformationiface.CloudFormationAPI
}

func (fc *fakeClient) ListStacks(input *cf.ListStacksInput) (*cf.ListStacksOutput, error) {
	if *input.StackStatusFilter[0] == "UPDATE_COMPLETE" {
		return &cf.ListStacksOutput{
			StackSummaries: []*cf.StackSummary{
				&cf.StackSummary{
					StackName:   aws.String("test2"),
					StackStatus: aws.String("UPDATE_COMPLETE"),
				},
			},
		}, nil

	}

	if *input.StackStatusFilter[0] == "ERROR" {
		return nil, errors.New("error")
	}

	return &cf.ListStacksOutput{
		StackSummaries: []*cf.StackSummary{
			&cf.StackSummary{
				StackName:   aws.String("test1"),
				StackStatus: aws.String("DELETE_COMPLETE"),
			},
		},
	}, nil
}

func TestListStacks(t *testing.T) {
	testData := []map[string]string{
		nil,
		map[string]string{"filter": "UPDATE_COMPLETE"},
		map[string]string{"filter": "ERROR"},
	}

	stack := NewStack(&fakeClient{})
	for _, td := range testData {
		err := stack.ListStacks("", td["filter"])
		if td["filter"] == "ERROR" {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
