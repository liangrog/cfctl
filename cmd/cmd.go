package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// cfctl config file
	cfgFileName      = ".cfctl"
	cfgFileExtension = ".yaml"
)

var (
	cfctlShort = i18n.T("cfctl manages stacks' lifecycle")

	cfctlLong = templates.LongDesc(i18n.T(`
		cfctl manages AWS CloudFormation stacks' lifecycle.

		For more information, please visit: https://github.com/liangrog/cfctl`))
)

var cfgFile string

var Cmds = NewCmdCfctl()

// Root cmd
func NewCmdCfctl() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "cfctl",
		Short: cfctlShort,
		Long:  cfctlLong,
	}

	return cmds
}

// Load config file and register persistent flags
func init() {
	cobra.OnInitialize(initConfig)

	Cmds.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is $HOME/%s%s)", cfgFileName, cfgFileExtension))
	Cmds.PersistentFlags().StringP(CMD_ROOT_OUTPUT, "o", "json", "output type. Default to json. Use 'yaml' for yaml output.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home := utils.HomeDir()

		fileName := fmt.Sprintf("%s%s", cfgFileName, cfgFileExtension)
		viper.AddConfigPath(home)
		viper.SetConfigName(fileName)

		// Create config file if doesn't exist
		abs := filepath.Join(home, fileName)
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			if _, err := os.Create(abs); err != nil {
				fmt.Println("Failed to create cfctl config file:", err)
				os.Exit(1)
			}
		}
	}

	//viper.SetConfigType("yaml")
	/*if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}*/
}

func Execute() {
	if err := Cmds.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}
