package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Register sub commands
func init() {
	cmd := getCmdTemplateValidate()
	addFlagsTemplateValidate(cmd)

	CmdTemplate.AddCommand(cmd)
}

func addFlagsTemplateValidate(cmd *cobra.Command) {
	cmd.Flags().BoolP(CMD_TEMPLATE_VALIDATE_RECURSIVE, "r", true, "Recursively validate templates for given path")
}

// cmd: validate
func getCmdTemplateValidate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate Cloudformation template. Example 'cfctl template validate [template path or s3 url]'",
		Long: `Validate Cloudformation template by given path.
This command can be run recursively by using '-r' option`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(utils.MsgFormat("One template path or url is required", utils.MessageTypeError))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			recursive, _ := cmd.Flags().GetBool(CMD_TEMPLATE_VALIDATE_RECURSIVE)

			err := templateValidate(
				cmd.Flags().Lookup(CMD_ROOT_OUTPUT).Value.String(),
				args,
				recursive,
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

type validOut struct {
	file string
	err  error
}

// Validate template
func templateValidate(format string, paths []string, recursive bool) error {
	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	for _, path := range paths {
		if ok, _ := utils.IsDir(path); ok {
			files, err := utils.FindFiles(path, recursive)
			if err != nil {
				return err
			}

			result := make(chan validOut, 10)
			for _, file := range files {
				go func(file string, res chan<- validOut) {
					res <- validOut{
						file: file,
						err:  tplValid(stack, file),
					}
				}(file, result)
			}

			for j := 1; j <= len(files); j++ {
				jobRes := <-result
				if jobRes.err != nil {
					return errors.New(fmt.Sprintf("%s %s", jobRes.file, jobRes.err))
				}
			}
		} else {
			// If only to validate a file or url
			err := tplValid(stack, path)

			if err != nil {
				return errors.New(fmt.Sprintf("%s %s", path, err))
			}
		}
	}

	utils.Print(utils.FormatType(format), "No error found")

	return nil
}

// Validate a simple template
func tplValid(stack *ctlaws.Stack, path string) error {
	var tplByte []byte
	var url string
	var err error

	if utils.IsUrlRegexp(path) {
		url = path
	} else {
		tplByte, err = ioutil.ReadFile(path)
		if err != nil {
			return err
		}
	}

	_, err = stack.ValidateTemplate(tplByte, url)

	return err
}
