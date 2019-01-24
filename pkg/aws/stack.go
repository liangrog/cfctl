package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/liangrog/cfctl/pkg/utils"
)

const (
	maxTemplateLength = 51200
)

// Stack struct.
// Provide API testing stub
type Stack struct {
	Client cloudformationiface.CloudFormationAPI
}

// Stack constructor
func NewStack(cfapi cloudformationiface.CloudFormationAPI) *Stack {
	return &Stack{Client: cfapi}
}

// List all stacks. Aggregate all pages and output only one array
func (s *Stack) ListStacks(format string, statusFilter ...string) ([]*cf.StackSummary, error) {
	var status []*string
	var nextToken *string
	var stackSummary []*cf.StackSummary

	if len(statusFilter) > 0 {
		status = aws.StringSlice(statusFilter)
	}

	for {
		input := &cf.ListStacksInput{
			NextToken:         nextToken,
			StackStatusFilter: status,
		}

		output, err := s.Client.ListStacks(input)
		if err != nil {
			return stackSummary, err
		}

		// aggregate all summaries
		//
		// Note: We don't expect a large
		// number of returns as it's limited
		// to the memory upper bound
		if len(output.StackSummaries) > 0 {
			stackSummary = append(stackSummary, output.StackSummaries...)
		}

		if output.NextToken == nil {
			break
		}

		nextToken = output.NextToken
	}

	return stackSummary, nil
}

// Validate cloudformation template
//
// url must be in AWS s3 URL. See https://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/#ValidateTemplateInput
//
func (s *Stack) ValidateTemplate(tpl []byte, url string) (*cf.ValidateTemplateOutput, error) {
	var input *cf.ValidateTemplateInput
	var output *cf.ValidateTemplateOutput

	// Must have one valide
	if len(tpl) == 0 && len(url) == 0 {
		return output, errors.New(utils.MsgFormat("Missing cloudformation template or template URLs", utils.MessageTypeError))
	}

	// If template string is given
	if len(tpl) > 0 {
		if len(tpl) > maxTemplateLength {
			return output, errors.New(utils.MsgFormat(fmt.Sprintf("Exceeded maximum template size of %d bytes", maxTemplateLength), utils.MessageTypeError))
		}

		input = &cf.ValidateTemplateInput{
			TemplateBody: aws.String(string(tpl)),
		}
	}

	// If only urls are given
	if len(tpl) == 0 && len(url) > 0 {
		input = &cf.ValidateTemplateInput{
			TemplateURL: aws.String(url),
		}

	}

	return s.Client.ValidateTemplate(input)
}

// Convert tags from map to Tag slice
func (s *Stack) TagSlice(tags map[string]string) []*cf.Tag {
	var t []*cf.Tag
	for k, v := range tags {
		t = append(t, new(cf.Tag).SetKey(k).SetValue(v))
	}

	return t
}

// Convert params from map to Parameter slice
func (s *Stack) ParamSlice(params map[string]string) []*cf.Parameter {
	var p []*cf.Parameter
	for k, v := range params {
		p = append(p, new(cf.Parameter).SetParameterKey(k).SetParameterValue(v))
	}

	return p
}

// Create a stack
func (s *Stack) CreateStack(name string, params map[string]string, tags map[string]string, tpl []byte, url string) (*cf.CreateStackOutput, error) {
	var stackOutput *cf.CreateStackOutput

	// Validate template
	Valid, err := s.ValidateTemplate(tpl, url)
	if err != nil {
		return stackOutput, err
	}

	input := new(cf.CreateStackInput).
		SetStackName(name).
		SetParameters(s.ParamSlice(params)).
		SetCapabilities(Valid.Capabilities).
		SetTags(s.TagSlice(tags))

	// Template
	if len(tpl) > 0 {
		input.SetTemplateBody(string(tpl))
	} else {
		input.SetTemplateURL(url)
	}

	return s.Client.CreateStack(input)
}

// Delete a stack
func (s *Stack) DeleteStack(stackName string, retainResc ...string) (*cf.DeleteStackOutput, error) {
	input := new(cf.DeleteStackInput).
		SetStackName(stackName)

	if len(retainResc) > 0 {
		input.SetRetainResources(aws.StringSlice(retainResc))
	}

	return s.Client.DeleteStack(input)
}

// Describe a stack by a given name
func (s *Stack) DescribeStack(stackName string) ([]*cf.Stack, error) {
	if len(stackName) <= 0 {
		return nil, errors.New(utils.MsgFormat("Missing stack name.", utils.MessageTypeError))
	}

	input := new(cf.DescribeStacksInput).SetStackName(stackName)
	out, err := s.Client.DescribeStacks(input)
	if err != nil {
		return nil, err
	}

	return out.Stacks, nil
}

// If stack exists
func (s *Stack) Exist(stackName string) bool {
	stacks, err := s.DescribeStack(stackName)

	if err != nil || len(stacks) != 1 {
		return false
	}

	return true
}

// Describe all stacks
func (s *Stack) DescribeStacks() ([]*cf.Stack, error) {
	var out []*cf.Stack
	var nextToken *string

	for {
		input := &cf.DescribeStacksInput{
			NextToken: nextToken,
		}

		o, err := s.Client.DescribeStacks(input)
		if err != nil {
			return out, err
		} else if len(o.Stacks) > 0 {
			out = append(out, o.Stacks...)
		}

		if o.NextToken == nil {
			break
		}

		nextToken = o.NextToken
	}

	return out, nil
}

// Get stack event
func (s *Stack) DescribeStackEvents(stackName string) ([]*cf.StackEvent, error) {
	var events []*cf.StackEvent
	var nextToken *string

	for {
		input := &cf.DescribeStackEventsInput{
			NextToken: nextToken,
			StackName: aws.String(stackName),
		}

		out, err := s.Client.DescribeStackEvents(input)
		if err != nil {
			return nil, err
		}

		events = append(events, out.StackEvents...)

		if out.NextToken == nil {
			break
		}

		nextToken = out.NextToken
	}

	return events, nil
}

// Kick off a stack drift detection process. Returns a
// detection process Id to be used for status query
func (s *Stack) DetectStackDrift(stackName string, resourceIds ...string) (string, error) {
	var detectionId string

	if len(stackName) == 0 {
		return detectionId, errors.New(utils.MsgFormat("Missing stack name.", utils.MessageTypeError))
	}

	input := new(cf.DetectStackDriftInput).
		SetStackName(stackName)

	if len(resourceIds) > 0 {
		input.SetLogicalResourceIds(aws.StringSlice(resourceIds))
	}

	output, err := s.Client.DetectStackDrift(input)
	if err != nil {
		return detectionId, err
	} else {
		detectionId = *output.StackDriftDetectionId
	}

	return detectionId, nil
}

// Detailing the stack drift at resources
func (s *Stack) DescribeStackResourceDrifts(stackName string, status ...string) ([]*cf.StackResourceDrift, error) {
	if len(stackName) == 0 {
		return nil, errors.New(utils.MsgFormat("Missing stack name.", utils.MessageTypeError))
	}

	input := new(cf.DescribeStackResourceDriftsInput).
		SetStackName(stackName)

	if len(status) > 0 {
		input.SetStackResourceDriftStatusFilters(aws.StringSlice(status))
	}

	output, err := s.Client.DescribeStackResourceDrifts(input)
	if err != nil {
		return nil, err
	}

	return output.StackResourceDrifts, nil
}

// Get current drift detection process status
func (s *Stack) DescribeStackDriftDetectionStatus(id string) (*cf.DescribeStackDriftDetectionStatusOutput, error) {
	return s.Client.DescribeStackDriftDetectionStatus(new(cf.DescribeStackDriftDetectionStatusInput).SetStackDriftDetectionId(id))
}
