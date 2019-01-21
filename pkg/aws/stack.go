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

// List all stacks and print to stdout
func (s *Stack) CmdListStacks(format string, statusFilter ...string) error {
	stackSummary, err := s.ListStacks(format, statusFilter...)
	if err != nil {
		return err
	}

	if err := utils.Print(format, stackSummary); err != nil {
		return err
	}

	return nil
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
		return output, errors.New(utils.MsgFormat("Missing cloudformation template or template URLs"))
	}

	// If template string is given
	if len(tpl) > 0 {
		if len(tpl) > maxTemplateLength {
			return output, errors.New(utils.MsgFormat(fmt.Sprintf("Exceeded maximum template size of %d bytes", maxTemplateLength)))
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
