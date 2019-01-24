package utils

import (
	"os"
	"testing"

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
