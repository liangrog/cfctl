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
	ENV_VAULT_PASSWORD = "ANSIBLE_VAULT_PASSWORD"
	CMD_VAULT_PASSWORD = "vault-password"
)

// Register sub commands
func init() {
	cmd := getCmdStackDeploy()
	addFlagsStackDeploy(cmd)

	CmdStack.AddCommand(cmd)
}

func addFlagsStackDeploy(cmd *cobra.Command) {
	cmd.Flags().BoolP("dry-run", "", false, "Validate the templates and parse the parameters but not creating the stack(s).")
	cmd.Flags().BoolP("quiet", "q", false, "Don't print out stack creation process such as events.")
	cmd.Flags().String("file", "", "Name of the stack configuration file (Default is './stacks.yaml')")
	cmd.Flags().String("set", "", "Set additional variables as key=value")
	cmd.Flags().String("env", "", "Set the environment you want to run")
	cmd.Flags().String("stack", "", "Specify stack(s)")
	cmd.Flags().String(CMD_VAULT_PASSWORD, "", "Ansible vault password, alternatively the password can be provided in environment variable ANSIBLE_VAULT_PASSWORD")
}

// cmd: create
func getCmdStackDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy stack(s)",
		Long:  `deploy CloudFormation stack(s)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			quiet, _ := cmd.Flags().GetBool("quiet")

			err := deployStacks(
				cmd.Flags().Lookup("set").Value.String(),
				cmd.Flags().Lookup("file").Value.String(),
				cmd.Flags().Lookup("env").Value.String(),
				cmd.Flags().Lookup("stack").Value.String(),
				GetVaultPass(cmd.Flags().Lookup(CMD_VAULT_PASSWORD).Value.String()),
				dryRun,
				quiet,
			)

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Get vault password from either cmd or env
func GetVaultPass(pass string) string {
	if len(pass) == 0 {
		pass = os.Getenv(ENV_VAULT_PASSWORD)
	}

	return pass
}

// Deploy stacks
func deployStacks(vars, f, env, named, vaultPass string, dry, quiet bool) error {
	var err error

	// Load deploy configuration file
	dc, err := conf.NewDeployConfig(f)
	if err != nil {
		return err
	}

	opts := map[string]interface{}{
		"quiet": quiet,
	}

	// Create S3 bucket if it doesn't exist
	cfs3 := ctlaws.NewS3(s3.New(ctlaws.AWSSess))
	if exist, err := cfs3.IfBucketExist(dc.S3Bucket); err != nil {
		return err
	} else if !exist {
		utils.CmdPrint(
			opts,
			utils.FormatCmd,
			utils.MsgFormat(fmt.Sprintf("Bucket %s doesn't exist. It will be created.", dc.S3Bucket), utils.MessageTypeInfo),
		)

		if _, err := cfs3.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(dc.S3Bucket)}); err != nil {
			return err
		}
	}

	// Retrieve the list of stacks
	var cmdsl []string
	if len(named) > 0 {
		cmdsl = strings.Split(named, " ")
	}

	sl, err := dc.GetStackList(cmdsl)
	if err != nil {
		return err
	}

	// Load key-value from env folders
	kv, err := LoadEnvValues(vaultPass, dc, env)
	if err != nil {
		return err
	}

	// Check if stacks are cyclic
	isCyclic, sorted, err := IfCircular(dc, sl, kv)
	if err != nil {
		return err
	} else if isCyclic {
		return errors.New("The stacks configuration contains circular dependency.")
	}

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	for _, v := range sorted {
		var stc *conf.StackConfig
		// Don't process if it's not in given stack list as it
		// may contains stacks from other references.
		stc, ok := sl[v.Value.Id()]
		if !ok {
			continue
		}

		// Get Parameters
		paramTpl, err := ioutil.ReadFile(dc.GetParamPath(stc.Param))
		if err != nil {
			return err
		}

		paramTpl, err = utils.GetCleanYamlBytes(paramTpl)
		if err != nil {
			return err
		}

		// Parse parameter template
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

		var waiterType string
		if stack.Exist(stc.Name) {
			_, err = stack.UpdateStack(stc.Name, params, stc.Tags, dat, "")
			waiterType = "update"
		} else {
			_, err = stack.CreateStack(stc.Name, params, stc.Tags, dat, "")
			waiterType = "create"
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

// Load specific environmennt values
func LoadEnvValues(vaultPass string, dc *conf.DeployConfig, env string) (map[string]string, error) {
	var values, envValues map[string]string

	// Load key-values from "default" folder
	if yes, err := utils.IsDir(dc.GetEnvDirPath("default")); err == nil && yes {
		values, err = conf.LoadValues(dc.GetEnvDirPath("default"), vaultPass)
		if err != nil {
			return nil, err
		}
	}

	// Load key-vaules from specified env folder
	if yes, err := utils.IsDir(dc.GetEnvDirPath(env)); len(env) > 0 && err == nil && yes {
		envValues, err = conf.LoadValues(dc.GetEnvDirPath(env), vaultPass)
		if err != nil {
			return nil, err
		}
	}

	// Combine key-values by overriding default values with env values
	for k, v := range envValues {
		values[k] = v
	}

	return values, nil
}

// Build a graph an check if there is any cyclic condition
func IfCircular(dc *conf.DeployConfig, sc map[string]*conf.StackConfig, kv map[string]string) (bool, []*gp.Vertice, error) {
	// Closure
	newVertice := func(name string) *gp.Vertice {
		vertice := gp.NewVertice(name, gl.NewEdgeStore())
		vertice.Value.SetId(name)
		return vertice
	}

	// Create a graph
	g := gp.NewGraph(gp.DIRECTED, gl.NewVerticeStore())

	for _, c := range sc {
		// Create a node for current stack
		vertice := g.GetVerticeById(c.Name)
		if vertice == nil {
			vertice = newVertice(c.Name)
		}

		g.AddVertice(vertice)

		// Search for dependent stacks
		content, err := ioutil.ReadFile(dc.GetParamPath(c.Param))
		if err != nil {
			return false, nil, err
		}

		content, err = utils.GetCleanYamlBytes(content)
		if err != nil {
			return false, nil, err
		}

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

			// Add new edge to vertice node
			edge := gp.NewEdge(p, dv, gp.FROM)
			vertice.AddEdge(edge)
		}

		// Update graph and amend the missing edges
		g.UpdateVertice(vertice)
	}

	return gs.Kahn(g)
}
