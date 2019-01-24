package cmd

import (
	"errors"
	"fmt"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdStackGet()
	addFlagsStackGet(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackGet(cmd *cobra.Command) {
	cmd.Flags().String("name", "n", "Get stack's details for given stack name")
}

// cmd: get
func getCmdStackGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get stack(s) details",
		Long: `Get all stacks details by default. If a stack name given, 
only return the detail for that stack`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := stackGet(
				cmd.Flags().Lookup("output").Value.String(),
				cmd.Flags().Lookup("name").Value.String(),
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Get stacks
func stackGet(format, stackName string) error {
	var out []*cf.Stack

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	// If stack name given
	if len(stackName) > 0 {
		if !stack.Exist(stackName) {
			return errors.New(utils.MsgFormat(fmt.Sprintf("Failed to find stack %s", stackName), utils.MessageTypeError))
		}

		if out, err := stack.DescribeStack(stackName); err != nil {
			return err
		} else {
			if err := utils.Print(format, out[0]); err != nil {
				return err
			}
		}

		return nil
	}

	out, err := stack.DescribeStacks()
	if err != nil {
		return err
	}

	if err := utils.Print(format, out); err != nil {
		return err
	}

	return nil
}
