package utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHomeDir(t *testing.T) {
	homeDir := "/home/cfctl"
	os.Setenv("HOME", homeDir)
	assert.Equal(t, homeDir, HomeDir())
}

func TestIsDir(t *testing.T) {
	ok, err := IsDir("./fsys_test.go")
	assert.False(t, ok)
	assert.NoError(t, err)

	ok, err = IsDir("../utils")
	assert.True(t, ok)
	assert.NoError(t, err)

}

func TestIsUrl(t *testing.T) {
	assert.True(t, IsUrl("https://google.com.au"))
	assert.True(t, IsUrl("/google/com/au"))
}

func TestIsUrlRegexp(t *testing.T) {
	assert.True(t, IsUrlRegexp("https://google.com.au"))
	assert.False(t, IsUrlRegexp("google.com.au"))
	assert.False(t, IsUrlRegexp("/google/com/au"))
}

func TestFindFiles(t *testing.T) {
	list, err := FindFiles("../../", true)
	assert.NoError(t, err)
	assert.True(t, len(list) > 0)
}

func TestScanFiles(t *testing.T) {
	unixTime := time.Now().Unix()
	tmpDir, err := ioutil.TempDir("", "test"+strconv.FormatInt(unixTime, 10))
	assert.NoError(t, err)

	_, err = ioutil.TempFile(tmpDir, "a")
	assert.NoError(t, err)

	_, err = ioutil.TempFile(tmpDir, "b")
	assert.NoError(t, err)

	// level 1 dir
	tmpDir2 := path.Join(tmpDir, "level2")
	err = os.Mkdir(tmpDir2, 0777)
	assert.NoError(t, err)

	_, err = ioutil.TempFile(tmpDir2, "c")
	assert.NoError(t, err)

	call := func(lvl int) ([]string, error) {
		var res []string
		done := make(chan bool)
		defer close(done)

		paths, errc := ScanFiles(tmpDir, done, lvl)
		for path := range paths {
			res = append(res, path)
		}

		err := <-errc
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	res, err := call(0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))

	res, err = call(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

	err = os.RemoveAll(tmpDir)
	assert.NoError(t, err)
}

func TestLoadYaml(t *testing.T) {
	unixTime := time.Now().Unix()
	tmpDir, err := ioutil.TempDir("", "test"+strconv.FormatInt(unixTime, 10))
	assert.NoError(t, err)

	_, err = ioutil.TempFile(tmpDir, "pkg-utils-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	assert.NoError(t, ioutil.WriteFile(path.Join(tmpDir, "pkg-utils-test"), []byte("test: blah\n#bascc"), 0644))

	c, err := LoadYaml(path.Join(tmpDir, "pkg-utils-test"))
	assert.NoError(t, err)

	scanner := bufio.NewScanner(strings.NewReader(string(c)))
	scanner.Split(bufio.ScanLines)

	// Count the lines.
	count := 0
	for scanner.Scan() {
		count++
	}
	// The comments in the file should be stripped out.
	assert.Equal(t, 1, count)
}
