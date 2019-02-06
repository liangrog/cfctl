package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

var CmdVault = getCmdVault()

const (
	ENV_PASSWORD          = "CFCTL_VAULT_PASSWORD"
	ENV_PASSWORD_FILE     = "CFCTL_VAULT_PASSWORD_FILE"
	DEFAULT_PASSWORD_FILE = ".cfctl_vault_password"
)

// Register sub commands
func init() {
	addFlagsVault(CmdVault)
	Cmds.AddCommand(CmdVault)
}

// cmd: vault
func getCmdVault() *cobra.Command {
	return &cobra.Command{
		Use:   "vault",
		Short: "Cipher with ansible-vault like",
		Long:  `Encrypt and decrypt following ansible-vault spec`,
	}
}

func addFlagsVault(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("password", "p", "", "Password for encryption or decryption")
	cmd.PersistentFlags().StringP("password-file", "", "", "Password file for encryption or decryption")
}

// getPassword function will attempt to locate passwrod
// via three ways in order:
// 1. CLI option --password and --password-file
// 2. Environment CFCTL_VAULT_PASSWORD or CFCTL_VAULT_PASSWORD_FILE
// 3. Default passwrod file in ~/.cfctl_vault_password
//
// Multiple passwords are seperated by ","
func GetPasswords(pass, passFile string) ([]string, error) {
	fileToSlice := func(path string) ([]string, error) {
		text, err := ioutil.ReadFile(path)
		if len(text) > 0 && err == nil {
			return strings.Split(string(text), ","), nil
		}

		return nil, err
	}

	// If password option provided
	if len(pass) > 0 {
		return strings.Split(pass, ","), nil
	}

	// If password-file option provide
	if len(passFile) > 0 {
		return fileToSlice(passFile)
	}

	// If env CFCTL_VAULT_PASSWORD provided
	if pass := os.Getenv(ENV_PASSWORD); len(pass) > 0 {
		return strings.Split(pass, ","), nil
	}

	// if env CFCTL_VAULT_PASSWORD_FILE provided
	if file := os.Getenv(ENV_PASSWORD_FILE); len(file) > 0 {
		return fileToSlice(file)
	}

	// Check default password file
	defaultFile := path.Join(utils.HomeDir(), DEFAULT_PASSWORD_FILE)
	if _, err := os.Stat(defaultFile); err == nil {
		return fileToSlice(defaultFile)
	}

	// Prompt password if all failed
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Password: ")
	if scanner.Scan() {
		passwords := scanner.Text()
		if len(passwords) > 0 {
			return strings.Split(passwords, ","), nil
		}
	}

	return nil, errors.New("Password is empty.")
}
