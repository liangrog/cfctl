package aws

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws/session"
)

func TestAWSSess(t *testing.T) {
	assert.IsType(t, new(session.Session), AWSSess)
}

func TestGetHttpClient(t *testing.T) {
	// no proxy
	assert.IsType(t, new(http.Client), GetHttpClient())

	// With proxy
	proxy := "https://127.0.0.2:3128"
	os.Setenv(ENV_HTTPS_PROXY, proxy)
	//Make sure proxy condition is triggered
	assert.Equal(t, proxy, getHTTPSProxy())
	assert.IsType(t, new(http.Client), GetHttpClient())
}

type httpsEnvSets struct {
	cap_val  string
	low_val  string
	expected string
}

func TestGetHTTPSProxy(t *testing.T) {
	testData := []httpsEnvSets{
		{},
		{cap_val: "https://127.0.0.1:3128", expected: "https://127.0.0.1:3128"},
		{low_val: "https://127.0.0.1:3128", expected: "https://127.0.0.1:3128"},
		{cap_val: "https://127.0.0.1:3128", low_val: "https://127.0.0.2:3128", expected: "https://127.0.0.2:3128"},
	}

	for _, d := range testData {
		os.Setenv(ENV_HTTPS_PROXY, d.cap_val)
		os.Setenv(strings.ToLower(ENV_HTTPS_PROXY), d.low_val)
		assert.Equal(t, d.expected, getHTTPSProxy())
	}
}
