package conf

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// filter function type
type stackFilter func(sc *StackConfig, s string) bool

// Filter function and value to use for
type filterConfig struct {
	// Needle for the filter
	value string

	// filter function
	funct stackFilter
}

// Filter by name
func nameFilter(sc *StackConfig, names string) bool {
	nameList := strings.Split(names, ",")
	for _, name := range nameList {
		if name == sc.Name {
			return true
		}
	}

	return false
}

// Tag filter
func tagFilter(sc *StackConfig, tags string) bool {
	tagList := strings.Split(tags, ",")
	c := len(tagList)
	for _, tag := range tagList {
		tagAttr := strings.Split(tag, "=")
		for k, v := range sc.Tags {
			if k == tagAttr[0] && v == tagAttr[1] {
				c--
			}
		}
	}

	// If all tags are matched
	if c == 0 {
		return true
	}

	return false
}

// Return all stack configs
func allFilter(sc *StackConfig, all string) bool {
	return true
}

// Get filter list
func getFilters(f map[string]string) []filterConfig {
	var sf []filterConfig

	// If no filter given, provide allfiler
	if f == nil || len(f) == 0 {
		sf = append(sf, filterConfig{funct: allFilter})
		return sf
	}

	for k, v := range f {
		switch k {
		case "name":
			sf = append(sf, filterConfig{value: v, funct: nameFilter})
		case "tag":
			sf = append(sf, filterConfig{value: v, funct: tagFilter})
		}
	}

	return sf
}

// Find stack config for given list
func (dc *DeployConfig) GetStackList(f map[string]string) map[string]*StackConfig {
	result := make(map[string]*StackConfig)

	filters := getFilters(f)
	// Get the full list in the right format
	for _, sc := range dc.Stacks {
		result[sc.Name] = sc
	}

	for k, sc := range result {
		for _, f := range filters {
			// If found not meet the given filter,
			// remove from the result
			if !f.funct(sc, f.value) {
				delete(result, k)
			}
		}
	}

	return result
}
