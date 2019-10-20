package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var (
	completionShort = i18n.T("Output shell completion code for the specified shell (bash or zsh).")

	completionLong = templates.LongDesc(i18n.T(`
		Output shell completion code for the specified shell (bash or zsh). 
		The shell code must be evaluated to provide interactive completion of cfctl commands.
		This can be done by sourcing it from the .bash_profile.

		Detailed instructions on how to do this are available here:
		https://github.com/liangrog/cfctl#enabling-shell-autocompletion

		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2
		`))

	completionExample = templates.Examples(i18n.T(`
		# Installing bash completion on macOS using homebrew

		## If running Bash 3.2 included with macOS
		brew install bash-completion

		## or, if running Bash 4.1+
		brew install bash-completion@2

		## You need add the completion to your completion directory
		cfctl completion bash > $(brew --prefix)/etc/bash_completion.d/cfctl


		# Installing bash completion on Linux

		## If bash-completion is not installed on Linux, please install the 'bash-completion' package
		## via your distribution's package manager. Take RedHat/Centos for example:
		yum install -y bash-completion bash-completion-extras

		## Load the cfctl completion code for bash into the current shell
		source <(cfctl completion bash)
		  
		## Write bash completion code to a file and source if from .bash_profile
		cfctl completion bash > ~/.cfctl/completion.bash.inc
		printf "
		# cfctl shell completion
		source '$HOME/.cfctl/completion.bash.inc'
		" >> $HOME/.bash_profile
		source $HOME/.bash_profile

		# Load the cfctl completion code for zsh[1] into the current shell
		source <(cfctl completion zsh)
		# Set the cfctl completion code for zsh[1] to autoload on startup
		 cfctl completion zsh > "${fpath[1]}/_cfctl"`))
)

// Register sub commands
func init() {
	cmdCompletion := getCmdCompletion()

	Cmds.AddCommand(cmdCompletion)
}

// cmd: stack
func getCmdCompletion() *cobra.Command {
	return &cobra.Command{
		Use:     "completion",
		Short:   completionShort,
		Long:    completionLong,
		Example: fmt.Sprintf(completionExample),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || len(args) > 1 {
				return errors.New("Need to provide shell type (bash or zch) in command argument")
			}

			if args[0] != "bash" && args[0] != "zsh" {
				return errors.New("Only 'bash' or 'zsh' allowed in the argument")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := completionGen(args[0])

			silenceUsageOnError(cmd, err)

			return err
		},
	}
}

// Generate shell scripts based on the shell type
// for auto completion.
func completionGen(shellType string) error {
	switch shellType {
	case "zsh":
		return Cmds.GenZshCompletion(os.Stdout)
	case "bash":
		return Cmds.GenBashCompletion(os.Stdout)
	}

	return nil
}
