package conf

import (
	"io/ioutil"
	"sort"
	"sync"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/vault"
	"gopkg.in/yaml.v2"
)

// Result for value read
type valueResult struct {
	path   string
	values map[string]string
	err    error
}

// Processing given files into key-values
func processValue(passwords []string, paths <-chan string, out chan<- *valueResult, done <-chan bool) {
	var v *valueResult
	var tv map[string]string

	for p := range paths {
		dat, err := ioutil.ReadFile(p)
		if err != nil {
			v = &valueResult{err: err}
		}
		// If it's vault encrypted file, decrypt it
		if vault.HasVaultHeader(dat) {
			decrypted := false
			for _, pass := range passwords {
				if dat, err = vault.Decrypt(pass, dat); err == nil {
					// If found one password that works.
					decrypted = true
					break
				}
			}

			// If there is a problem, don't continue
			if !decrypted {
				out <- &valueResult{err: err}
				continue
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
func LoadValues(root string, passwords []string) (map[string]string, error) {
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
			processValue(passwords, paths, c, done)
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
