package cmd

import (
	"errors"
	_ "fmt"
	_ "io/ioutil"

	_ "github.com/liangrog/cfctl/pkg/template/parser"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdTemplateParse()
	addFlagsTemplateParse(cmd)

	CmdTemplate.AddCommand(cmd)
}

func addFlagsTemplateParse(cmd *cobra.Command) {
	cmd.Flags().String("root-path", "", "Folder path to where stack.yaml is located")
}

// cmd: parse
func getCmdTemplateParse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse files that contain template functions",
		Long:  `Parse files that contain template functions`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(utils.MsgFormat("One file path for parsing is required", utils.MessageTypeError))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := parseTemplate(
				cmd.Flags().Lookup("root-path").Value.String(),
				args,
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

func parseTemplate(rootPath string, tpl []string) error {
	/*
		dat, err := ioutil.ReadFile(tpl[0])
		if err != nil {
			return err
		}

		bucket := "https://test.s3.amazonaws.com"

		data := make(map[string]string)
		data["VpcId"] = "1234"
		data["SecurityGroupTplUrl"] = "{{ tpl \"asg/default.yaml\" }}"
		out, err := parser.Parse(string(dat), bucket, data)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", out)
	*/
	return nil
}
