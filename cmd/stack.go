package cmd

import (
	"github.com/spf13/cobra"
)

var CmdStack = getCmdStack()

// Register sub commands
func init() {
	Cmds.AddCommand(CmdStack)
}

// cmd: stack
func getCmdStack() *cobra.Command {
	return &cobra.Command{
		Use:   "stack",
		Short: "manage stacks",
		Long:  `All actions that manages stack(s)`,
	}
}
