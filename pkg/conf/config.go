package conf

import (
	"io/ioutil"
	"sort"
	"sync"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/vault"
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

	// Values for the template
	ValueDir string `yaml:"valueDir,omitempty"`

	// Stacks config
	Stacks []*StackConfig `yaml:"stacks"`
}

type StackConfig struct {
	Name       string `yaml:"Name"`
	Template   string `yaml:"template"`
	Parameters string `yaml:"parameters,omitempty"`
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

	dc := new(DeployConfig)
	if err := yaml.Unmarshal(data, dc); err != nil {
		return nil, err
	}

	return dc, nil
}

// Result for value read
type valueResult struct {
	path   string
	values map[string]string
	err    error
}

// Processing given files into key-values
func processValue(password string, paths <-chan string, out chan<- *valueResult, done <-chan bool) {
	var v *valueResult
	var tv map[string]string

	for p := range paths {
		dat, err := ioutil.ReadFile(p)
		if err != nil {
			v = &valueResult{err: err}
		}
		// If it's vault encrypted file, decrypt it
		if vault.HasVaultHeader(dat) {
			dat, err = vault.Decrypt(password, dat)
			if err != nil {
				v = &valueResult{err: err}
			}
		}

		// Keep overriding the same map
		err = yaml.Unmarshal([]byte(dat), &tv)
		if err != nil {
			v = &valueResult{err: err}
		} else {
			v = &valueResult{
				path:   p,
				values: tv,
			}
		}

		select {
		case out <- v:
		case <-done:
			return
		}
	}
}

// Load Values from value directories.
// The order of value override is always following
// the lexical order of the file names
func LoadValues(root string, password string) (map[string]string, error) {
	// Find all files in paths
	done := make(chan bool)
	defer close(done)

	// Start file scanning
	paths, errc := utils.ScanFiles(root, done, 0)

	// Start 20 workers
	var wg sync.WaitGroup
	numProc := 10
	wg.Add(numProc)

	c := make(chan *valueResult)
	for i := 0; i < numProc; i++ {
		go func() {
			processValue(password, paths, c, done)
			wg.Done()
		}()
	}

	// Close result when all workers
	go func() {
		wg.Wait()
		close(c)
	}()

	// Processing output

	m := make(map[string]map[string]string)
	for vr := range c {
		if vr.err != nil {
			return nil, vr.err
		}
		m[vr.path] = vr.values
	}

	// Check whether the file scan failed.
	if err := <-errc; err != nil {
		return nil, err
	}

	// Get keys
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	// Merge values
	var values map[string]string
	sort.Strings(keys)
	for _, k := range keys {
		values = MergeValues(values, m[k])
	}

	return values, nil
}

// Merge two maps
func MergeValues(target, replace map[string]string) map[string]string {
	if target == nil {
		return replace
	}

	for k, v := range replace {
		target[k] = v
	}

	return target
}

func GetDependencyTree() {
}
