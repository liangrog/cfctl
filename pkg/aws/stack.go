package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/liangrog/cfctl/pkg/utils"
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
func (s *Stack) ListStacks(format string, statusFilter ...string) error {
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
			return err
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

	err := utils.Print(format, stackSummary)
	if err != nil {
		return err
	}

	return nil
}

// Validate cloudformation template
func (s *Stack) ValidateTemplate(tplStr []byte) error {
}
