package cmd

import (
	"errors"
	"fmt"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var (
	stackGetShort = i18n.T("Get one or more stack details")

	stackGetLong = templates.LongDesc(i18n.T(`
		Get all stacks details by default. 
		If stack names given, only return the detail for those stacks`))

	stackGetExample = templates.Examples(i18n.T(`
		# Get all stack details in config file backend-infra.yaml
		$ cfctl stack get --file backend-infra.yaml

		# Get a specific stack 'stack-a' detail
		$ cfctl stack get --name stack-a 

		# Get multiple stack details
		$ cfctl stack get --name stack-a,stack-b

		# Get stack details with tag Name=frontend
		$ cfctl stack get --tags Name=frontend`))
)

// Register sub commands
func init() {
	cmd := getCmdStackGet()
	addFlagsStackGet(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackGet(cmd *cobra.Command) {
	cmd.Flags().String(CMD_STACK_GET_NAME, "", "Get stack's details for given stack name. Multiple stack names can be given and seperated by comma, e.g 'stack-a,stack-b'")
}

// cmd: get
func getCmdStackGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   stackGetShort,
		Long:    stackGetLong,
		Example: fmt.Sprintf(stackGetExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := stackGet(
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_FILE).Value.String(),
				cmd.Flags().Lookup(CMD_ROOT_OUTPUT).Value.String(),
				cmd.Flags().Lookup(CMD_STACK_GET_NAME).Value.String(),
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_TAGS).Value.String(),
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Get stacks
func stackGet(f, format, stackNames, tags string) error {
	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	// Load deploy configuration file.
	dc, err := conf.NewDeployConfig(f)
	if err != nil {
		return err
	}

	// Retrieve the list of stacks and apply filters.
	filters := make(map[string]string)
	if len(stackNames) > 0 {
		filters["name"] = stackNames
	}

	if len(tags) > 0 {
		filters["tag"] = tags
	}

	sl := dc.GetStackList(filters)

	if len(sl) == 0 {
		return errors.New("No stack found.")
	}

	// If stack name given
	var errMsg []string
	for k, _ := range sl {
		if !stack.Exist(k) {
			errMsg = append(errMsg, utils.MsgFormat(fmt.Sprintf("Failed to find stack %s\n", k), utils.MessageTypeError))
			continue
		}

		if out, err := stack.DescribeStack(k); err != nil {
			return err
		} else {
			if err := utils.Print(utils.FormatType(format), out); err != nil {
				return err
			}
		}
	}

	// Print out error message
	if len(errMsg) > 0 {
		for _, msg := range errMsg {
			utils.StdoutWarn(msg)
		}
	}

	return nil
}
