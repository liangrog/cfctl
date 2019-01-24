package cmd

import (
	"github.com/spf13/cobra"
)

var CmdTemplate = getCmdTemplate()

// Register sub commands
func init() {
	Cmds.AddCommand(CmdTemplate)
}

// cmd: template
func getCmdTemplate() *cobra.Command {
	return &cobra.Command{
		Use:   "template",
		Short: "Cloudformation template utitlity",
		Long:  `Utilities for managing Cloudformation templates`,
	}
}
