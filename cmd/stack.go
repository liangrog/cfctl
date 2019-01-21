package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/spf13/cobra"
)

var stack *ctlaws.Stack

// Register sub commands
func init() {
	stack = ctlaws.NewStack(cloudformation.New(ctlaws.AWSSess))

	listStack := getCmdListStacks()
	Cmds.AddCommand(listStack)
	addFlagsListStacks(listStack)
}

// cmd: list-stacks
func getCmdListStacks() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list-stacks",
		Short: "List all stacks",
		Long:  `List all existing CloudFormation stacks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var status []string
			if sfilter := cmd.Flags().Lookup("status").Value.String(); len(sfilter) > 0 {
				status = strings.Split(sfilter, ",")
			}

			return stack.CmdListStacks(cmd.Flags().Lookup("output").Value.String(), status...)
		},
	}

	return cmd
}

func addFlagsListStacks(cmd *cobra.Command) {
	cmd.Flags().String("status", "", "cloudformation status filter, multiple values seperate by ','. Allowed values 'REVIEW_IN_PROGRESS, CREATE_FAILED, UPDATE_ROLLBACK_FAILED, UPDATE_ROLLBACK_IN_PROGRESS, CREATE_IN_PROGRESS, UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS, ROLLBACK_IN_PROGRESS, DELETE_COMPLETE, UPDATE_COMPLETE, UPDATE_IN_PROGRESS, DELETE_FAILED, DELETE_IN_PROGRESS, ROLLBACK_COMPLETE, ROLLBACK_FAILED, UPDATE_COMPLETE_CLEANUP_IN_PROGRESS, CREATE_COMPLETE, UPDATE_ROLLBACK_COMPLETE'")
}
