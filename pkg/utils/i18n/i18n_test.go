package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslation(t *testing.T) {
	sample := "Dragon slayer"

	assert.Equal(t, sample, T(sample))
}
