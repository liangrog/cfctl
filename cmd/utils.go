package cmd

import (
	"github.com/spf13/cobra"
)

// Silence usage output and downstream
// error output when returns error
// using RunE in command
func silenceUsageOnError(cmd *cobra.Command, err error) {
	if err != nil {
		//		cmd.SilenceUsage = true
	}
}
