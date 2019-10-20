package cmd

import (
	"fmt"
	"strings"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var (
	stackListShort = i18n.T("List all stacks")

	stackListLong = templates.LongDesc(i18n.T(`List all existing CloudFormation stacks`))

	stackListExample = templates.Examples(i18n.T(``))
)

// Register sub commands
func init() {
	cmd := getCmdStackList()
	addFlagsStackList(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackList(cmd *cobra.Command) {
	cmd.Flags().String(CMD_STACK_LIST_STATUS, "s", "cloudformation status filter, multiple values seperate by ','. Allowed values 'REVIEW_IN_PROGRESS, CREATE_FAILED, UPDATE_ROLLBACK_FAILED, UPDATE_ROLLBACK_IN_PROGRESS, CREATE_IN_PROGRESS, UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS, ROLLBACK_IN_PROGRESS, DELETE_COMPLETE, UPDATE_COMPLETE, UPDATE_IN_PROGRESS, DELETE_FAILED, DELETE_IN_PROGRESS, ROLLBACK_COMPLETE, ROLLBACK_FAILED, UPDATE_COMPLETE_CLEANUP_IN_PROGRESS, CREATE_COMPLETE, UPDATE_ROLLBACK_COMPLETE'")
}

// cmd: list
func getCmdStackList() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Short:   stackListShort,
		Long:    stackListLong,
		Example: fmt.Sprintf(stackListExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			var status []string

			if cmd.Flags().Changed(CMD_STACK_LIST_STATUS) {
				status = strings.Split(
					cmd.Flags().Lookup(CMD_STACK_LIST_STATUS).Value.String(),
					",",
				)
			}

			err := listStacks(cmd.Flags().Lookup(CMD_ROOT_OUTPUT).Value.String(), status...)

			silenceUsageOnError(cmd, err)

			return err
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

	if err := utils.Print(utils.FormatType(format), stackSummary); err != nil {
		return err
	}

	return nil
}
