package cmd

import (
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var CmdStack = getCmdStack()

var (
	stackShort = i18n.T("Commands for managing stack lifecycle.")

	stackLong = templates.LongDesc(i18n.T(`Manage stack lifecycle such as creation, update and deletion`))
)

// Register sub commands
func init() {
	Cmds.AddCommand(CmdStack)

	Cmds.PersistentFlags().StringP(CMD_STACK_DEPLOY_FILE, "f", "", "Alternative stack configuration file (Default is './stacks.yaml')")
	Cmds.PersistentFlags().StringP(CMD_STACK_DEPLOY_TAGS, "", "", "Only run stacks that match the specified tags in the form of 'tag=value'. Multiple tags can be given seperated by comma, e.g. 'tag1=value1,tag2=value2'. If stack names being provided at the argument at the same time, it will use both for filtering.")
}

// cmd: stack
func getCmdStack() *cobra.Command {
	return &cobra.Command{
		Use:   "stack",
		Short: stackShort,
		Long:  stackLong,
	}
}
