package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/liangrog/cfctl/pkg/utils/editor"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/liangrog/vault"
	"github.com/spf13/cobra"
)

const (
	// cfctl editor env
	EDITOR_ENV = "CFCTL_EDITOR"
)

var (
	vaultEditShort = i18n.T("Edit vault file in a seperate terminal")

	vaultEditLong = templates.LongDesc(i18n.T(`
		Edit an encrypted vault file by openning a seperate
		terminal using default editor. 'CFCTL_VAULT_PASSWORD'
		and 'CFCTL_VAULT_PASSWORD_FILE' environment variables 
		can be used to replace '--vault-password' and 
		'--vault-password-file' flags.`))

	vaultEditExample = templates.Examples(i18n.T(`
		# Edit a vault file named 'secrect.yaml'
		cfctl vault edit secret.yaml

		# Edit a vault file 'secret.yaml' with a password argument
		cfctl vault edit secret.yaml --vault-password Password1

		# Edit a vault file 'secret.yaml' providing a vault password file in $HOME/vault-password-file
		cfctl vault edit secret.yaml --vault-password-file $HOME/vault-password-file`))
)

// Register sub commands
func init() {
	cmd := getCmdVaultEdit()

	CmdVault.AddCommand(cmd)
}

// cmd: encrypt
func getCmdVaultEdit() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   vaultEditShort,
		Long:    vaultEditLong,
		Example: fmt.Sprintf(vaultEditExample),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || len(args) > 1 {
				return errors.New("Need to provide one filename only in command argument")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := editVaultFile(
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD).Value.String(),
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD_FILE).Value.String(),
				args,
			)

			silenceUsageOnError(cmd, err)

			return err

		},
	}

	return cmd
}

// Edit encrypted vault file in a seperate terminal using default editor
func editVaultFile(password, passwordFile string, files []string) error {
	passwords, err := GetPasswords(password, passwordFile, false, false)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		return err
	}

	// Try every given password
	var originalPlainText []byte
	var foundPass string

	decrypted := false
	for _, p := range passwords {
		originalPlainText, err = vault.Decrypt(p, data)
		if err == nil {
			decrypted = true
			foundPass = p
			break
		}
	}

	if !decrypted {
		return errors.New(fmt.Sprintf("Failed to decrypt %s using all given password", files[0]))
	}

	buf := bytes.NewBuffer(originalPlainText)

	edit := editor.NewDefaultEditor([]string{EDITOR_ENV})
	editedPlainText, _, err := edit.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(files[0])), filepath.Ext(files[0]), buf)
	if err != nil {
		return err
	}

	encryptedText, err := vault.Encrypt(editedPlainText, foundPass)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(files[0], encryptedText, 0644); err != nil {
		return err
	}

	return nil
}
