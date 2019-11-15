package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var CmdVault = getCmdVault()

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
	cmd.PersistentFlags().StringP(CMD_VAULT_PASSWORD, "", "", "vault password for encryption or decryption")
	cmd.PersistentFlags().StringP(CMD_VAULT_PASSWORD_FILE, "", "", "file that contains vault passwords for encryption or decryption")
}

// getPassword function will attempt to locate passwrod
// via three ways in order:
// 1. CLI option --password and --password-file
// 2. Environment CFCTL_VAULT_PASSWORD or CFCTL_VAULT_PASSWORD_FILE
// 3. Default passwrod file in ~/.cfctl_vault_password
//
// Multiple passwords are seperated by ","
func GetPasswords(pass, passFile string, noPrompt, allowEmpty bool) ([]string, error) {
	fileToSlice := func(path string) ([]string, error) {
		var result []string
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 {
				result = append(result, (strings.Split(line, ","))...)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return result, nil
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
	if pass := os.Getenv(ENV_VAULT_PASSWORD); len(pass) > 0 {
		return strings.Split(pass, ","), nil
	}

	// if env CFCTL_VAULT_PASSWORD_FILE provided
	if file := os.Getenv(ENV_VAULT_PASSWORD_FILE); len(file) > 0 {
		return fileToSlice(file)
	}

	// Check default password file
	defaultFile := path.Join(utils.HomeDir(), DEFAULT_VAULT_PASSWORD_FILE)
	if _, err := os.Stat(defaultFile); err == nil {
		return fileToSlice(defaultFile)
	}

	// If empty password is allowed, return
	if allowEmpty {
		return []string{""}, nil
	}

	// Prompt password if all failed
	if !noPrompt {
		fmt.Print("Password: ")
		if passwords, err := terminal.ReadPassword(int(syscall.Stdin)); err == nil {
			// Make sure cursor start a new line
			fmt.Print("\n")
			if len(passwords) > 0 {
				return strings.Split(strings.TrimSpace(string(passwords)), ","), nil
			}
		}
	}

	return nil, errors.New("Password is empty.")
}
