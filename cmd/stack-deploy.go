package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/template/parser"
	"github.com/liangrog/cfctl/pkg/utils"
	gl "github.com/liangrog/ds/graph/list"
	gp "github.com/liangrog/ds/graph/parts"
	gs "github.com/liangrog/ds/graph/sort"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	// Environment variable for vault password.
	ENV_VAULT_PASSWORD = "ANSIBLE_VAULT_PASSWORD"

	// Command line flag for vault password.
	CMD_VAULT_PASSWORD = "vault-password"

	// Command line flag for stacks.
	CMD_STACK_DEPLOY_STACK = "stack"

	// Command line flag for configuration file.
	CMD_STACK_DEPLOY_FILE = "file"

	// Command line flag for dry run.
	CMD_STACK_DEPLOY_DRY_RUN = "dry-run"

	// Command line flag for envoirnment folder.
	CMD_STACK_DEPLOY_ENV = "env"

	// Default environment folder name
	STACK_DEPLOY_ENV_DEFAULT_FOLDER = "default"
)

// Register sub commands.
func init() {
	cmd := getCmdStackDeploy()
	addFlagsStackDeploy(cmd)

	CmdStack.AddCommand(cmd)
}

// Add flags to stack deploy command.
func addFlagsStackDeploy(cmd *cobra.Command) {
	cmd.Flags().BoolP(CMD_STACK_DEPLOY_DRY_RUN, "", false, "Validate the templates and parse the parameters but not creating the stacks")
	cmd.Flags().String(CMD_STACK_DEPLOY_FILE, "", "Alternative stack configuration file (Default is './stacks.yaml')")
	cmd.Flags().String(CMD_STACK_DEPLOY_ENV, "", "Set enviornment folder you want to load values from")
	cmd.Flags().String(CMD_STACK_DEPLOY_STACK, "", "Specify what stacks to run. If multiple stacks, use comma delimiter. For example: stackA,stackB")
	cmd.Flags().String(CMD_VAULT_PASSWORD, "", "Ansible vault password, alternatively the password can be provided in environment variable ANSIBLE_VAULT_PASSWORD")
}

// cmd: stack deploy.
func getCmdStackDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy stacks",
		Long:  `deploy CloudFormation stacks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool(CMD_STACK_DEPLOY_DRY_RUN)

			err := deployStacks(
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_FILE).Value.String(),
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_ENV).Value.String(),
				cmd.Flags().Lookup(CMD_STACK_DEPLOY_STACK).Value.String(),
				getVaultPass(cmd.Flags().Lookup(CMD_VAULT_PASSWORD).Value.String()),
				dryRun,
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Get vault password from either command line flag or envirnment variable.
func getVaultPass(pass string) string {
	if len(pass) == 0 {
		pass = os.Getenv(ENV_VAULT_PASSWORD)
	}

	return pass
}

// Load key-value from a givenn environmennt folder.
func loadEnvValues(vaultPass string, dc *conf.DeployConfig, envFolder string) (map[string]string, error) {
	var values, envValues map[string]string

	// Load key-values from "default" folder.
	if yes, err := utils.IsDir(dc.GetEnvDirPath(STACK_DEPLOY_ENV_DEFAULT_FOLDER)); err == nil && yes {
		values, err = conf.LoadValues(dc.GetEnvDirPath(STACK_DEPLOY_ENV_DEFAULT_FOLDER), vaultPass)
		if err != nil {
			return nil, err
		}
	}

	// Load key-vaules from specified env folder.
	if yes, err := utils.IsDir(dc.GetEnvDirPath(envFolder)); len(envFolder) > 0 && err == nil && yes {
		envValues, err = conf.LoadValues(dc.GetEnvDirPath(envFolder), vaultPass)
		if err != nil {
			return nil, err
		}
	}

	// Combine key-values by overriding default values with env values.
	for k, v := range envValues {
		values[k] = v
	}

	return values, nil
}

// Check if there is any cyclic dependency condition in stacks being
// deployed and return a ordered list, providing priority to parent
// stacks by building a graph and doing a Kahn sort.
func ifCircularStacks(dc *conf.DeployConfig, sc map[string]*conf.StackConfig, kv map[string]string) (bool, []*gp.Vertice, error) {
	// Anonymous function for creating vertice .
	newVertice := func(name string) *gp.Vertice {
		vertice := gp.NewVertice(name, gl.NewEdgeStore())
		// Set vertice Id to stack name so
		// we can search by it.
		vertice.Value.SetId(name)
		return vertice
	}

	// Create a directed graph.
	g := gp.NewGraph(gp.DIRECTED, gl.NewVerticeStore())

	// Loop through the stacks to create vertices.
	for _, c := range sc {
		// Create a node for current stack.
		vertice := g.GetVerticeById(c.Name)
		if vertice == nil {
			vertice = newVertice(c.Name)
		}

		g.AddVertice(vertice)

		// Load parameter file
		content, err := utils.LoadYaml(dc.GetParamPath(c.Param))
		if err != nil {
			return false, nil, err
		}

		// Search for dependent stacks.
		dep, err := parser.SearchDependancy(string(content), kv)
		if err != nil {
			return false, nil, err
		}

		// Add dependent stacks as nodes.
		for _, p := range dep {
			dv := g.GetVerticeById(p)

			if dv == nil {
				dv = newVertice(p)
				g.AddVertice(dv)
			}

			// Add new edge to child node.
			edge := gp.NewEdge(p, dv, gp.FROM)
			vertice.AddEdge(edge)
		}

		// Update graph and amend the missing edges.
		g.UpdateVertice(vertice)
	}

	return gs.Kahn(g)
}

// Deploy stacks.
func deployStacks(f, env, named, vaultPass string, dry bool) error {
	var err error

	// Load deploy configuration file.
	dc, err := conf.NewDeployConfig(f)
	if err != nil {
		return err
	}

	// Create S3 bucket if it doesn't exist.
	cfs3 := ctlaws.NewS3(s3.New(ctlaws.AWSSess))
	if exist, err := cfs3.IfBucketExist(dc.S3Bucket); err != nil {
		return err
	} else if !exist {
		utils.InfoPrint(fmt.Sprintf("[ warning ] s3 bucket %s doesn't exist. It will be created.", dc.S3Bucket), utils.MessageTypeInfo)

		if _, err := cfs3.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(dc.S3Bucket)}); err != nil {
			return err
		}
	}

	utils.InfoPrint(fmt.Sprintf("[ info ] found s3 bucket %s", dc.S3Bucket))

	// Retrieve the list of stacks from comma delimited string.
	var cmdsl []string
	if len(named) > 0 {
		cmdsl = strings.Split(named, ",")
	}

	sl, err := dc.GetStackList(cmdsl)
	if err != nil {
		return err
	}

	// Load key-value from env folder.
	kv, err := loadEnvValues(vaultPass, dc, env)
	if err != nil {
		return err
	}

	// Check if stacks are cyclic and sort it.
	isCyclic, sorted, err := ifCircularStacks(dc, sl, kv)
	if err != nil {
		return err
	} else if isCyclic {
		return errors.New("The stacks configuration contains circular dependency.")
	}

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	for _, v := range sorted {
		// Don't process if it's not in given stack list as it
		// may contains stacks from other references such via
		// tpl function.
		stc, ok := sl[v.Value.Id()]
		if !ok {
			continue
		}

		// Get Parameters.
		paramTpl, err := utils.LoadYaml(dc.GetParamPath(stc.Param))
		if err != nil {
			return err
		}

		// Line seperator for each stack
		fmt.Println("")

		// Parse parameter template.
		paramBytes, err := parser.Parse(string(paramTpl), kv, dc)
		if err != nil {
			return err
		}

		params := make(map[string]string)
		if err := yaml.Unmarshal(paramBytes, &params); err != nil {
			return err
		}

		dat, err := ioutil.ReadFile(dc.GetTplPath(stc.Tpl))
		if err != nil {
			return err
		}

		// Dry run
		if dry {
			if _, err := stack.ValidateTemplate(dat, ""); err != nil {
				return err
			} else {
				utils.InfoPrint(
					fmt.Sprintf(
						"[ stack | validate ] %s\t%s",
						stc.Name,
						"ok",
					),
				)
			}

			continue
		}

		// Create or update the stack.
		var waiterType string
		if stack.Exist(stc.Name) {
			_, err = stack.UpdateStack(stc.Name, params, stc.Tags, dat, "")
			waiterType = ctlaws.StackWaiterTypeUpdate
		} else {
			_, err = stack.CreateStack(stc.Name, params, stc.Tags, dat, "")
			waiterType = ctlaws.StackWaiterTypeCreate
		}

		if err != nil {
			return err
		}

		if err := stack.PollStackEvents(stc.Name, waiterType); err != nil {
			return err
		}
	}

	return nil
}
