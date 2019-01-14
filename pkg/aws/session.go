// AWS session facilities. Allow using https proxy for requests
package aws

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	// Environment variable names
	// https proxy. This allows using
	// https proxy other than bash's http proxy
	ENV_HTTPS_PROXY_CAP = "CF_HTTPS_PROXY"
	ENV_HTTPS_PROXY_LOW = "cf_https_proxy"
)

var (
	AWSSess = session.Must(
		session.NewSession(
			&aws.Config{
				HTTPClient: GetHttpClient(),
			},
		),
	)
)

// Get http client for aws calls
func GetHttpClient() *http.Client {
	// Setup tool specific https proxy if available
	// This allows cfctl to use proxy other than
	// standard 'https_proxy'
	if proxyStr := getHTTPSProxy(); len(proxyStr) > 0 {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			glog.Fatalf("Error parsing proxy URL %s", proxyStr)
		}

		transport := http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{},
		}

		return &http.Client{
			Transport: &transport,
		}
	}

	return &http.Client{}
}

// Get https proxy string from
// environment variable
func getHTTPSProxy() string {
	var proxyStr string

	if len(os.Getenv(ENV_HTTPS_PROXY_CAP)) > 0 {
		proxyStr = os.Getenv(ENV_HTTPS_PROXY_CAP)
	}

	// Override capital variable if set
	if len(os.Getenv(ENV_HTTPS_PROXY_LOW)) > 0 {
		proxyStr = os.Getenv(ENV_HTTPS_PROXY_LOW)
	}

	return proxyStr
}
