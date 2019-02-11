package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/vault"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdVaultEncrypt()

	CmdVault.AddCommand(cmd)
}

// cmd: encrypt
func getCmdVaultEncrypt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt given content",
		Long:  `Encrypt given content following ansible-vault spec`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New(utils.MsgFormat("Missing file name in command argument", utils.MessageTypeError))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := batchEncrypt(
				cmd.Flags().Lookup("password").Value.String(),
				cmd.Flags().Lookup("password-file").Value.String(),
				args,
			)

			silenceUsageOnError(cmd, err)

			return err

		},
	}

	return cmd
}

func batchEncrypt(pss, pssFile string, files []string) error {
	passwords, err := GetPasswords(pss, pssFile)
	if err != nil {
		return err
	}

	// Only one password allowed
	if len(passwords) > 1 {
		return errors.New("More than one passwords were given")
	}

	result := make(chan error, 10)
	for _, file := range files {
		go func(file, pass string, res chan<- error) {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				res <- err
				return
			}

			output, err := vault.Encrypt(data, pass)
			if err != nil {
				res <- err
				return
			}

			if err := ioutil.WriteFile(file, output, 0644); err != nil {
				res <- err
			}

			res <- nil
		}(file, passwords[0], result)
	}

	for j := 0; j < len(files); j++ {
		err := <-result
		if err != nil {
			if err := utils.Print("", files[j], err); err != nil {
				return err
			}
		}
	}

	return nil
}
