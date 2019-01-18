package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdoutStrFactory(t *testing.T) {
	assert.IsType(t, StdoutStrFactory("yaml"), StdoutStrFactory(""))
}

func TestPrint(t *testing.T) {
	assert.Error(t, Print(""))
	assert.NoError(t, Print("", "testing Print"))
}
