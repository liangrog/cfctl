package cmd

import (
	"strings"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdStackList()
	addFlagsStackList(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackList(cmd *cobra.Command) {
	cmd.Flags().String("status", "s", "cloudformation status filter, multiple values seperate by ','. Allowed values 'REVIEW_IN_PROGRESS, CREATE_FAILED, UPDATE_ROLLBACK_FAILED, UPDATE_ROLLBACK_IN_PROGRESS, CREATE_IN_PROGRESS, UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS, ROLLBACK_IN_PROGRESS, DELETE_COMPLETE, UPDATE_COMPLETE, UPDATE_IN_PROGRESS, DELETE_FAILED, DELETE_IN_PROGRESS, ROLLBACK_COMPLETE, ROLLBACK_FAILED, UPDATE_COMPLETE_CLEANUP_IN_PROGRESS, CREATE_COMPLETE, UPDATE_ROLLBACK_COMPLETE'")
}

// cmd: list
func getCmdStackList() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all stacks",
		Long:  `List all existing CloudFormation stacks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var status []string

			if sfilter := cmd.Flags().Lookup("status").Value.String(); len(sfilter) > 0 {
				status = strings.Split(sfilter, ",")
			}

			return listStacks(cmd.Flags().Lookup("output").Value.String(), status...)
		},
	}

	return cmd
}

// List all stacks and print to stdout
func listStacks(format string, statusFilter ...string) error {

	stackSummary, err := ctlaws.
		NewStack(cf.New(ctlaws.AWSSess)).
		ListStacks(format, statusFilter...)

	if err != nil {
		return err
	}

	if err := utils.Print(format, stackSummary); err != nil {
		return err
	}

	return nil
}
