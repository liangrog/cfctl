package utils

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
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
			if err == nil && !info.IsDir() {
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
