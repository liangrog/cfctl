package cmd

import (
	"errors"
	"fmt"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdStackDelete()
	addFlagsStackDelete(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackDelete(cmd *cobra.Command) {
	cmd.Flags().BoolP("all", "", false, "Delete all the stacks in the stack configuration file")
	cmd.Flags().String(CMD_STACK_DEPLOY_FILE, "", "Alternative stack configuration file (Default is './stacks.yaml')")
}

// cmd: delete
func getCmdStackDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete cloudformation stacks",
		Long:  `delete cloudformation stacks`,
		Args: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")

			if !all && len(args) < 1 {
				return errors.New("Minimum one stack name is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			all, _ := cmd.Flags().GetBool("all")
			err := stackDelete(
				args,
				all,
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_FILE).Value.String(),
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Delete stacks.
func stackDelete(stackNames []string, all bool, stackConf string) error {
	var err error

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	// If flag is set to all stacks, get
	// stacks from configuration file.
	if all {
		stackNames, err = getAllStacksFromConfig(stackConf)
		if err != nil {
			return err
		}
	}

	for _, sn := range stackNames {
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

// Get all stack names from configuration
// file and return in a slice.
func getAllStacksFromConfig(stackConf string) ([]string, error) {
	var list []string
	// Load deploy configuration file.
	dc, err := conf.NewDeployConfig(stackConf)
	if err != nil {
		return list, err
	}

	for _, sc := range dc.Stacks {
		list = append(list, sc.Name)
	}

	return list, nil
}
