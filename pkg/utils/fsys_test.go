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
