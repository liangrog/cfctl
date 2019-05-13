package utils

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// HomeDir returns the home directory for the current user
func HomeDir() string {
	if runtime.GOOS == "windows" {

		// First prefer the HOME environmental variable
		if home := os.Getenv("HOME"); len(home) > 0 {
			if _, err := os.Stat(home); err == nil {
				return home
			}
		}

		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDir := homeDrive + homePath
			if _, err := os.Stat(homeDir); err == nil {
				return homeDir
			}
		}

		if userProfile := os.Getenv("USERPROFILE"); len(userProfile) > 0 {
			if _, err := os.Stat(userProfile); err == nil {
				return userProfile
			}
		}
	}
	return os.Getenv("HOME")
}

// If given path a directory
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// If given string a url
// This validate follows RFC standard
func IsUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	return true
}

// Validate a url by using regexp
func IsUrlRegexp(urlStr string) bool {
	exp := `^(http[s]?:\/\/?){1}([0-9A-Za-z-\.@:%_\+~#=]+)+((\.[a-zA-Z]{2,3})+)(\/(.)*)?(\?(.)*)?`

	re := regexp.MustCompile(exp)
	return re.MatchString(urlStr)
}

// Search files in a directory
// A Simple implementation
func FindFiles(filePath string, recursive bool) ([]string, error) {
	var list []string

	// Anonymous function to append found json file to result
	add := func(path string, info os.FileInfo) {
		if !info.IsDir() {
			list = append(list, path)
		}
	}

	// If it's a recursive search, we'll walk
	if recursive {
		walk := func(path string, info os.FileInfo, err error) error {
			if err == nil && info.Mode().IsRegular() {
				add(path, info)
			}

			return nil
		}

		filepath.Walk(filePath, walk)
	} else {
		items, err := ioutil.ReadDir(filePath)
		if err != nil {
			return list, err
		}
		for _, info := range items {
			p := filepath.Join(filePath, info.Name())
			add(p, info)
		}
	}

	sort.Strings(list)

	return list, nil
}

// This function provides a fast result ready machanism
// via feeding the result immediately as soon as the first
// file found.
// By default it will scan recursively. You can define
// The level of the directory if desire.
// Given level 0 means recursively.
// Level 1 is only files in given root directory.
// Level 2 is files in root directory + any files in folders under root.
// Other levels are so on and so forth.
// Exist channel provides a machanism for exiting the scan early
// without creating memory leak.
func ScanFiles(root string, exit <-chan bool, level int) (<-chan string, <-chan error) {
	// Buffered channel, non-blocking
	files := make(chan string)
	e := make(chan error, 1)

	go func() {
		// Send signal to let listeners know all files have been scanned
		defer close(files)

		// Walk is following lexical order
		e <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// If not file, return
			if !info.Mode().IsRegular() {
				return nil
			}

			// Only search within level of directory limit
			if level > 0 {
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}

				parts := strings.Split(rel, "/")

				if len(parts) > level {
					return nil
				}
			}

			select {
			case files <- path:
			case <-exit:
				return errors.New("Recieved exit signal, existed")
			}

			return nil
		})
	}()

	return files, e
}

// Yaml cleansing such as remove comment in yaml file.
func GetCleanYamlBytes(input []byte) ([]byte, error) {
	t := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(input, &t); err != nil {
		return nil, err
	}

	return yaml.Marshal(&t)
}

// Load yaml file and return clean yaml bytes.
func LoadYaml(path string) ([]byte, error) {
	var result []byte
	// Read file.
	result, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}

	return GetCleanYamlBytes(result)
}

// Rewrite given path by returning from
// the first occurence of given string
func RewritePath(from, match string) string {
	p := strings.Split(from, "/")

	var index int
	for idx, v := range p {
		if v == match {
			index = idx
			break
		}
	}

	return path.Join(p[index:]...)
}
