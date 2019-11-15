package aws

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/stretchr/testify/assert"
)

var stack = NewStack(&stackFakeClient{})

// mock client
type stackFakeClient struct {
	cloudformationiface.CloudFormationAPI
}

func (fc *stackFakeClient) ListStacks(input *cf.ListStacksInput) (*cf.ListStacksOutput, error) {
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

func (fc *stackFakeClient) ValidateTemplate(input *cf.ValidateTemplateInput) (*cf.ValidateTemplateOutput, error) {
	return new(cf.ValidateTemplateOutput).SetDescription("testing"), nil
}

func (fc *stackFakeClient) CreateStack(input *cf.CreateStackInput) (*cf.CreateStackOutput, error) {
	return new(cf.CreateStackOutput).SetStackId("testing"), nil
}

func (fc *stackFakeClient) DeleteStack(input *cf.DeleteStackInput) (*cf.DeleteStackOutput, error) {
	return new(cf.DeleteStackOutput), nil
}

func (fc *stackFakeClient) DescribeStacks(input *cf.DescribeStacksInput) (*cf.DescribeStacksOutput, error) {
	var stacks []*cf.Stack

	sampleStack := new(cf.Stack).
		SetStackName("test").
		SetStackStatus(cf.StackStatusCreateComplete)

	stacks = append(stacks, sampleStack)

	return &cf.DescribeStacksOutput{
		Stacks: stacks,
	}, nil
}

func (fc *stackFakeClient) WaitUntilStackCreateComplete(input *cf.DescribeStacksInput) error {
	return nil
}

func (fc *stackFakeClient) DescribeStackEvents(input *cf.DescribeStackEventsInput) (*cf.DescribeStackEventsOutput, error) {
	var events []*cf.StackEvent

	e := new(cf.StackEvent).
		SetEventId("test-event").
		SetStackId("test-stack-id").
		SetStackName("test").
		SetTimestamp(time.Now()).
		SetLogicalResourceId("122345").
		SetResourceStatus("Complete")

	events = append(events, e)

	return &cf.DescribeStackEventsOutput{
		StackEvents: events,
	}, nil
}

func (fc *stackFakeClient) DetectStackDrift(input *cf.DetectStackDriftInput) (*cf.DetectStackDriftOutput, error) {
	return new(cf.DetectStackDriftOutput).SetStackDriftDetectionId("detect-id-123abc"), nil
}

func (fc *stackFakeClient) DescribeStackResourceDrifts(input *cf.DescribeStackResourceDriftsInput) (*cf.DescribeStackResourceDriftsOutput, error) {
	var drifts []*cf.StackResourceDrift

	d := new(cf.StackResourceDrift).
		SetStackId("test-abc").
		SetStackResourceDriftStatus(cf.StackDriftStatusDrifted)

	drifts = append(drifts, d)

	return &cf.DescribeStackResourceDriftsOutput{
		StackResourceDrifts: drifts,
	}, nil
}

func (fc *stackFakeClient) DescribeStackDriftDetectionStatus(input *cf.DescribeStackDriftDetectionStatusInput) (*cf.DescribeStackDriftDetectionStatusOutput, error) {
	return new(cf.DescribeStackDriftDetectionStatusOutput).
		SetDetectionStatus(cf.StackDriftDetectionStatusDetectionComplete).
		SetStackDriftDetectionId("abc-test").
		SetStackDriftStatus(cf.StackDriftStatusDrifted).
		SetStackId("test"), nil
}

func (fc *stackFakeClient) DescribeStackResources(input *cf.DescribeStackResourcesInput) (*cf.DescribeStackResourcesOutput, error) {
	var stackRes []*cf.StackResource

	sr := &cf.StackResource{
		LogicalResourceId: aws.String("abcd-1234"),
		ResourceStatus:    aws.String("CREATE_COMPLETE"),
		ResourceType:      aws.String("AWS::S3::Bucket"),
	}

	sr.SetTimestamp(time.Now())

	stackRes = append(stackRes, sr)

	return new(cf.DescribeStackResourcesOutput).
		SetStackResources(stackRes), nil
}

func TestTagSlice(t *testing.T) {
	data := map[string]string{
		"Name": "testing",
	}

	tags := NewStack(&stackFakeClient{}).TagSlice(data)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, "testing", *tags[0].Value)
}

func TestParamSlice(t *testing.T) {
	data := map[string]string{
		"S3Name": "testing",
	}

	params := NewStack(&stackFakeClient{}).ParamSlice(data)
	assert.Equal(t, 1, len(params))
	assert.Equal(t, "testing", *params[0].ParameterValue)
}

func TestListStacks(t *testing.T) {
	testData := []map[string]string{
		nil,
		map[string]string{"filter": "UPDATE_COMPLETE"},
		map[string]string{"filter": "ERROR"},
	}

	for _, td := range testData {
		stackSum, err := stack.ListStacks("", td["filter"])
		if td["filter"] == "ERROR" {
			assert.Error(t, err)
		} else if td["filter"] == "UPDATE_COMPLETE" {
			assert.NoError(t, err)
			assert.Equal(t, "test2", *stackSum[0].StackName)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, 1, len(stackSum))
		}
	}
}

func TestValidateTemplate(t *testing.T) {
	var tpl []byte
	var url string

	// test no params
	_, err := stack.ValidateTemplate(tpl, url)
	assert.Error(t, err)

	// test url
	_, err = stack.ValidateTemplate(tpl, "https://s3")
	assert.NoError(t, err)
}

func TestCreateStack(t *testing.T) {
	_, err := stack.CreateStack("testing", nil, nil, nil, "https://s3")
	assert.NoError(t, err)
}

func TestDeleteStack(t *testing.T) {
	_, err := stack.DeleteStack("testing")
	assert.NoError(t, err)
}

func TestDescribeStack(t *testing.T) {
	s, err := stack.DescribeStack("")
	assert.Error(t, err)

	s, err = stack.DescribeStack("test")
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestDescribeStacks(t *testing.T) {
	s, err := stack.DescribeStacks()
	assert.NoError(t, err)
	assert.True(t, len(s) > 0)
}

func TestDetectStackDrift(t *testing.T) {
	id, err := stack.DetectStackDrift("test")
	assert.NoError(t, err)
	assert.True(t, len(id) > 0)
}

func TestDescribeStackResourceDrifts(t *testing.T) {
	out, err := stack.DescribeStackResourceDrifts("test", cf.StackDriftStatusDrifted)
	assert.NoError(t, err)
	assert.True(t, len(out) > 0)
}

func TestDescribeStackDriftDetectionStatus(t *testing.T) {
	out, err := stack.DescribeStackDriftDetectionStatus("test")
	assert.NoError(t, err)
	assert.IsType(t, new(cf.DescribeStackDriftDetectionStatusOutput), out)
}

func TestPollStackEvents(t *testing.T) {
	assert.NoError(t, stack.PollStackEvents("test", StackWaiterTypeCreate))
}

func TestGetStackResources(t *testing.T) {
	out, err := stack.GetStackResources("test-stack")
	assert.NoError(t, err)
	assert.Equal(t, aws.StringValue(out[0].LogicalResourceId), "abcd-1234")
}
