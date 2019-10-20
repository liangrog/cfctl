package cmd

import (
	"errors"
	"fmt"
	"strings"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var (
	stackDeleteShort = i18n.T("Delete one or more stacks.")

	stackDeleteLong = templates.LongDesc(i18n.T(`Delete one or more stacks.`))

	stackDeleteExample = templates.Examples(i18n.T(`
		# Delete a stack with name 'stack-1'
		$ cfctl stack delete stack-1

		# Delete multiple stacks with name 'stack-1' and 'stack-2'
		$ cfctl stack delete stack-1 stack-2

		# Delete all stacks from a specific stack file
		$ cfctl stack delete --file stack-file.yaml --all
	
		# Delete stacks that have specific tag values
		$ cfctl sack delete --tags Name=stack-1,Type=frontend`))
)

// Register sub commands
func init() {
	cmd := getCmdStackDelete()
	addFlagsStackDelete(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackDelete(cmd *cobra.Command) {
	cmd.Flags().BoolP(CMD_STACK_DELETE_ALL, "", false, "Delete all the stacks in the stack configuration file")
}

// cmd: delete
func getCmdStackDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   stackDeleteShort,
		Long:    stackDeleteLong,
		Example: fmt.Sprintf(stackDeleteExample),
		Args: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool(CMD_STACK_DELETE_ALL)

			tags := cmd.Flags().Lookup(CMD_STACK_DEPLOY_TAGS).Value.String()

			if !all && len(tags) == 0 && len(args) < 1 {
				return errors.New(fmt.Sprintf("Please provide either stack name or using '--%s' or '--%s' flag", CMD_STACK_DELETE_ALL, CMD_STACK_DEPLOY_TAGS))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool(CMD_STACK_DELETE_ALL)
			err := stackDelete(
				args,
				all,
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_FILE).Value.String(),
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_TAGS).Value.String(),
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Delete stacks.
func stackDelete(stackNames []string, all bool, stackConf, tags string) error {
	var err error

	stacks := make(map[string]*conf.StackConfig)

	// Load deploy configuration file.
	dc, err := conf.NewDeployConfig(stackConf)
	if err != nil {
		return err
	}

	// If flag is set to all stacks, get
	// stacks from configuration file.
	if all {
		stacks = dc.GetStackList(nil)
	} else {
		filters := make(map[string]string)
		if len(tags) > 0 {
			filters["tag"] = tags
		}

		if len(stackNames) > 0 {
			filters["name"] = strings.Join(stackNames, ",")
		}

		stacks = dc.GetStackList(filters)
	}

	if len(stacks) == 0 {
		return errors.New(fmt.Sprintf("No stack found for given filters.\n"))
	}

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))
	for sn, _ := range stacks {
		fmt.Println("")

		// If stack name given
		if !stack.Exist(sn) {
			utils.StdoutError(fmt.Sprintf("Failed to find stack %s", sn))
			continue
		}

		_, err = stack.DeleteStack(sn)
		if err != nil {
			return err
		}

		if err := stack.PollStackEvents(sn, ctlaws.StackWaiterTypeDelete); err != nil {
			return err
		}
	}

	return nil
}
