package cmd

import (
	"github.com/spf13/cobra"
)

var CmdS3 = getCmdS3()

// Register sub commands
func init() {
	Cmds.AddCommand(CmdS3)
}

// cmd: s3
func getCmdS3() *cobra.Command {
	return &cobra.Command{
		Use:   "s3",
		Short: "s3 utilities",
		Long:  `Utilities for s3`,
	}
}
