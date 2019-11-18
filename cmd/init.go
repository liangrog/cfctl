package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/liangrog/cfctl/pkg/sample"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/spf13/cobra"
)

var (
	initShort = i18n.T("Initiate a sample repository.")

	initLong = templates.LongDesc(i18n.T(`
		Create a sample repository. It contains a simple folder structure
		and sample S3 bucket cloudformation and stack configuration so you
		can run a sample stack to test the tool.`))

	initExample = templates.Examples(i18n.T(`
		# Create a repo in /tmp with name cfctl-sample
		cfctl init --path /tmp --name cfctl-sample

		# Create sample stack 
		
		# You can either use -f option to specify the location of the
		# stack config file or
		cd /tmp/cfctl-sample
		cfctl stack deploy`))
)

// Register sub commands
func init() {
	cmd := getCmdInit()
	addFlagsInit(cmd)
	Cmds.AddCommand(cmd)
}

func addFlagsInit(cmd *cobra.Command) {
	cmd.Flags().String(CMD_INIT_PATH, "", "path where you want to create the repository.")
	cmd.Flags().String(CMD_INIT_NAME, "cfctl-sample", "name of the new repository. Default to 'cfctl'")
}

// cmd: init
func getCmdInit() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Short:   initShort,
		Long:    initLong,
		Example: fmt.Sprintf(initExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := sampleGen(
				cmd.Flags().Lookup(CMD_INIT_PATH).Value.String(),
				cmd.Flags().Lookup(CMD_INIT_NAME).Value.String(),
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}
}

// Generate sample repository
func sampleGen(p, name string) error {
	if len(name) == 0 {
		return errors.New("Missing sample repository name for init command")
	}

	dp := "./"
	if len(p) > 0 {
		dp = p
	}

	repo := path.Join(dp, name)

	// Create repository
	if err := os.MkdirAll(repo, 0755); err != nil {
		return err
	}

	// Create sub folders
	subDir := []string{
		"templates",
		"deploy",
		"deploy/sample",
		"deploy/sample/parameters",
		"deploy/sample/environments/default",
	}

	for _, d := range subDir {
		if err := os.MkdirAll(path.Join(repo, d), 0755); err != nil {
			return err
		}
	}

	// Create sample files
	files := map[string]string{
		"deploy/sample/stacks.yaml":                   sample.StackYaml,
		"deploy/sample/parameters/s3.yaml":            sample.SampleParam,
		"deploy/sample/environments/default/var.yaml": sample.EnvVars,
		"templates/s3-encrypted.yaml":                 sample.S3Template,
	}

	for k, v := range files {
		if err := writeSample(v, path.Join(repo, k)); err != nil {
			return err
		}
	}

	return nil
}

// Write sample file
func writeSample(s, file string) error {
	var f *os.File
	defer f.Close()

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(s)); err != nil {
		return err
	}

	return nil
}
