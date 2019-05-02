package conf

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/liangrog/cfctl/pkg/utils"
	"gopkg.in/yaml.v2"
)

const (
	// Default deployment package config file name
	DEFAULT_DEPLOY_CONFIG_FILE_NAME = "stacks.yaml"
)

// Deploy configuration
type DeployConfig struct {
	// Name of the s3 bucket for uploading template
	S3Bucket string `yaml:"s3Bucket"`

	// Template directory
	TemplateDir string `yaml:"templateDir"`

	// Environments directory
	EnvDir string `yaml:"envDir"`

	// Template directory
	ParamDir string `yaml:"paramDir"`

	// Stacks config
	Stacks []*StackConfig `yaml:"stacks"`

	// config file absolute path
	absPath string
}

// Stack configuration
type StackConfig struct {
	// Name of the stack
	Name string `yaml:"name"`

	// Template relative path
	Tpl string `yaml:"tpl"`

	// Parameter file relative path
	Param string `yaml:"param,omitempty"`

	Tags map[string]string `yaml:"tags,omitempty"`
}

// Load deploy config from file.
// If no file path given, default to lookup
// file "stacks.yaml" at current directory.
func NewDeployConfig(file string) (*DeployConfig, error) {
	if len(file) == 0 {
		file = DEFAULT_DEPLOY_CONFIG_FILE_NAME
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Allow using ENV function in configuration.
	// Parse environment variable.
	funcEnv := func(key string) string {
		return os.Getenv(key)
	}

	funcMap := template.FuncMap{"env": funcEnv}

	tmpl, err := template.New(uuid.New().String()).Funcs(funcMap).Parse(string(data))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, nil); err != nil {
		return nil, err
	}

	dc := new(DeployConfig)
	if err := yaml.Unmarshal(b.Bytes(), dc); err != nil {
		return nil, err
	}

	dc.absPath, err = filepath.Abs(filepath.Dir(file))
	if err != nil {
		return nil, err
	}

	if err := dc.Validate(); err != nil {
		return nil, err
	}

	return dc, nil
}

// Validate path configuration
func (dc *DeployConfig) Validate() error {
	var msg string

	if ok, err := utils.IsDir(path.Join(dc.absPath, dc.TemplateDir)); !ok || err != nil {
		msg = "There is a problem with templateDir in configuration."
	}

	if ok, err := utils.IsDir(path.Join(dc.absPath, dc.EnvDir)); !ok || err != nil {
		msg = "There is a problem with envDir in configuration."
	}

	if ok, err := utils.IsDir(path.Join(dc.absPath, dc.ParamDir)); !ok || err != nil {
		msg = "There is a problem with paramDir in configuration."
	}

	if len(msg) > 0 {
		return errors.New(utils.MsgFormat(msg, utils.MessageTypeError))
	}

	return nil
}

func (dc *DeployConfig) GetTplPath(n string) string {
	return path.Join(dc.absPath, dc.TemplateDir, n)
}

func (dc *DeployConfig) GetParamPath(n string) string {
	return path.Join(dc.absPath, dc.ParamDir, n)
}

func (dc *DeployConfig) GetEnvDirPath(n string) string {
	return path.Join(dc.absPath, dc.EnvDir, n)
}

// Return a stack config by its name
func (dc *DeployConfig) GetStackConfigByName(n string) *StackConfig {
	for _, sc := range dc.Stacks {
		if sc.Name == n {
			return sc
		}
	}

	return nil
}

// Find stack config for given list
func (dc *DeployConfig) GetStackList(l []string) (map[string]*StackConfig, error) {
	result := make(map[string]*StackConfig)
	if len(l) > 0 {
		for _, sc := range dc.Stacks {
			for _, s := range l {
				if sc.Name == s {
					result[s] = sc
				}
			}
		}

		if len(l) != len(result) {
			return nil, errors.New("There is mismatch between selected stacks and the configuration.")
		}

	} else {
		// Return full list if no selected stacks
		for _, sc := range dc.Stacks {
			result[sc.Name] = sc
		}
	}

	return result, nil
}
