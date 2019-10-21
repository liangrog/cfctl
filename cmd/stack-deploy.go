package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
	ctlaws "github.com/liangrog/cfctl/pkg/aws"
	"github.com/liangrog/cfctl/pkg/conf"
	"github.com/liangrog/cfctl/pkg/template/parser"
	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	gl "github.com/liangrog/ds/graph/list"
	gp "github.com/liangrog/ds/graph/parts"
	gs "github.com/liangrog/ds/graph/sort"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	stackDeployShort = i18n.T("Create or update one or more stacks")

	stackDeployLong = templates.LongDesc(i18n.T(`
		A single command that will create or update (if exists) one or more stacks
		depending given flags.`))

	stackDeployExample = templates.Examples(i18n.T(`
		# Deploy all stacks without using variable.
		$ cfctl stack deploy 

		# Deploy all stacks from a specific stack file
		$ cfctl stack deploy -f stack-file.yaml

		# Deploy particular stacks using variables from specific environment
		$ cfctl stack deploy --stack stack1,stack2 --env production

		# Deploy stacks using variables from specific environment that contains secrets and providing password file
		$ cfctl stack deploy --env production --vault-password-file path/to/password/file

		# Override environment values
		$ cfctl stack deploy --env production --vault-password-file path/to/password/file --vars name1=value1,name2=value2

		# Deploy stacks with specify tag values
		$ cfctl stack deploy --stack stack1,stack2 --tags Type=frontend

		# Output parameters only for all stacks
		$ cfctl stack deploy --env production --param-only`))
)

// Register sub commands.
func init() {
	cmd := getCmdStackDeploy()
	addFlagsStackDeploy(cmd)

	CmdStack.AddCommand(cmd)
}

// Add flags to stack deploy command.
func addFlagsStackDeploy(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(CMD_VAULT_PASSWORD, "", "", "Vault password for encryption or decryption")
	cmd.PersistentFlags().StringP(CMD_VAULT_PASSWORD_FILE, "", "", "File that contains vault passwords for encryption or decryption")

	cmd.Flags().BoolP(CMD_STACK_DEPLOY_DRY_RUN, "", false, "Validate the templates and parse the parameters but not creating the stacks")
	cmd.Flags().BoolP(CMD_STACK_DEPLOY_PARAM_ONLY, "", false, "Only parsing the parameter files")
	cmd.Flags().String(CMD_STACK_DEPLOY_ENV, "", "Set enviornment folder you want to load values from")
	cmd.Flags().String(CMD_STACK_DEPLOY_STACK, "", "Specify what stacks to run. If multiple stacks, use comma delimiter. For example: stackA,stackB")
	cmd.Flags().String(CMD_STACK_DEPLOY_VARS, "", "Specify variable override in the format of 'name=value'. If multiple , use comma delimiter.")
}

// cmd: stack deploy.
func getCmdStackDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   stackDeployShort,
		Long:    stackDeployLong,
		Example: fmt.Sprintf(stackDeployExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, _ := cmd.Flags().GetBool(CMD_STACK_DEPLOY_DRY_RUN)
			paramOnly, _ := cmd.Flags().GetBool(CMD_STACK_DEPLOY_PARAM_ONLY)
			passes, err := GetPasswords(
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD).Value.String(),
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD_FILE).Value.String(),
				false,
				true,
			)

			if err == nil {
				err = deployStacks(
					cmd.Flags().Lookup(CMD_STACK_DEPLOY_FILE).Value.String(),
					cmd.Flags().Lookup(CMD_STACK_DEPLOY_ENV).Value.String(),
					cmd.Flags().Lookup(CMD_STACK_DEPLOY_STACK).Value.String(),
					cmd.Flags().Lookup(CMD_STACK_DEPLOY_TAGS).Value.String(),
					passes,
					dryRun,
					paramOnly,
					cmd.Flags().Lookup(CMD_ROOT_OUTPUT).Value.String(),
					cmd.Flags().Lookup(CMD_STACK_DEPLOY_VARS).Value.String(),
				)
			}

			silenceUsageOnError(cmd, err)

			return err
		},
	}

	return cmd
}

// Load key-value from a givenn environmennt folder.
func loadEnvValues(vaultPass []string, dc *conf.DeployConfig, envFolder string) (map[string]string, error) {
	values := make(map[string]string)
	envValues := make(map[string]string)

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

		// Only add other depending vertices
		// and edges if parameters file exist.
		if len(c.Param) > 0 {
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
	}

	return gs.Kahn(g)
}

// Deploy stacks.
func deployStacks(f, env, named, tags string, vaultPass []string, dry, paramOnly bool, output string, vars string) error {
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
		utils.StdoutWarn(fmt.Sprintf("s3 bucket %s doesn't exist. It will be created.", dc.S3Bucket))

		if _, err := cfs3.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(dc.S3Bucket)}); err != nil {
			return err
		}
	} else {
		utils.StdoutInfo(fmt.Sprintf("found s3 bucket %s\n", dc.S3Bucket))
	}

	// Retrieve the list of stacks and apply filters.
	filters := make(map[string]string)

	if len(named) > 0 {
		filters["name"] = named
	}

	if len(tags) > 0 {
		filters["tag"] = tags
	}

	sl := dc.GetStackList(filters)

	// If no stack found, send a warning.
	if len(sl) == 0 {
		utils.StdoutWarn(fmt.Sprintf("No stack found for given filters. No further actions.\n"))
		return nil
	}

	// Load key-value from env folder.
	kv, err := loadEnvValues(vaultPass, dc, env)
	if err != nil {
		return err
	}

	// Load var override
	if len(vars) > 0 {
		varlist := strings.Split(vars, ",")
		for _, v := range varlist {
			vkv := strings.Split(v, "=")
			kv[vkv[0]] = vkv[1]
		}
	}

	// Check all stacks in the config file if it's cyclic
	fullList := dc.GetStackList(nil)
	isCyclic, _, err := ifCircularStacks(dc, fullList, kv)
	if err != nil {
		return err
	} else if isCyclic {
		return errors.New("The stack(s) in the stack list contains circular dependency.")
	}

	// Check if stacks are cyclic and sort it.
	_, sorted, _ := ifCircularStacks(dc, sl, kv)

	stack := ctlaws.NewStack(cf.New(ctlaws.AWSSess))

	for _, v := range sorted {
		// Don't process if it's not in given stack list as it
		// may contains stacks from other references such via
		// tpl function.
		stc, ok := sl[v.Value.Id()]
		if !ok {
			continue
		}

		// Line seperator for each stack
		fmt.Println("")

		// If there is parameters provided
		params := make(map[string]string)
		// If no parameters and only parsing parameters
		if len(stc.Param) <= 0 && paramOnly {
			continue
		}

		if len(stc.Param) > 0 {
			// Get Parameters.
			paramTpl, err := utils.LoadYaml(dc.GetParamPath(stc.Param))
			if err != nil {
				return err
			}

			// Parse parameter template.
			paramBytes, err := parser.Parse(string(paramTpl), kv, dc)
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(paramBytes, &params); err != nil {
				return err
			}

			// If only parsing parameters
			if paramOnly {
				if output == "yaml" {
					utils.InfoPrint("------")
					utils.InfoPrint(string(paramBytes))
				} else {
					pList := stack.ParamSlice(params)
					if pListJson, err := json.MarshalIndent(pList, "  ", "  "); err != nil {
						return err
					} else {
						utils.InfoPrint(string(pListJson))
					}
				}

				continue
			}

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
			if excludeErrorByMessage(err, stc.Name) {
				continue
			}
			return err
		}

		if err := stack.PollStackEvents(stc.Name, waiterType); err != nil {
			return err
		}
	}

	return nil
}

// Excluding some AWS errors.
func excludeErrorByMessage(err error, name string) bool {
	// No update error
	noUpdate := "No updates are to be performed"
	if strings.Contains(err.Error(), noUpdate) {
		utils.StdoutInfo(fmt.Sprintf("%s for %s\n", noUpdate, name))
		return true
	}

	return false
}
