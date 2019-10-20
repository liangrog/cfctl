package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	yamltext = `
---
s3Bucket: test
templateDir: {{ env "CF_TEST_TEMPLATE_DIR" }}
envDir: {{ env "CF_TEST_ENV_DIR" }}
paramDir: {{ env "CF_TEST_PARAM_DIR" }}
stacks:
  - name: stack-a
    tpl: stack-a.yaml
    tags:
      Name: stack-a
      App: test
  - name: stack-b
    tpl: stack-a.yaml
    tags:
      Name: stack-b
      App: test`
)

func setup(t *testing.T) (string, string) {
	// Create tmp folder
	unixTime := time.Now().Unix()
	tmpDir, err := ioutil.TempDir("", "test"+strconv.FormatInt(unixTime, 10))
	assert.NoError(t, err)

	templateDir := fmt.Sprintf("%s/templates", tmpDir)
	envDir := fmt.Sprintf("%s/env", tmpDir)
	paramDir := fmt.Sprintf("%s/param", tmpDir)

	// Create folders
	err = os.MkdirAll(templateDir, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(envDir, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(paramDir, os.ModePerm)
	assert.NoError(t, err)

	// Set folder env
	os.Setenv("CF_TEST_TEMPLATE_DIR", "templates")
	os.Setenv("CF_TEST_ENV_DIR", "env")
	os.Setenv("CF_TEST_PARAM_DIR", "param")

	stackFile, err := ioutil.TempFile(tmpDir, "stack.yaml.")
	assert.NoError(t, err)
	_, err = stackFile.Write([]byte(yamltext))
	assert.NoError(t, err)

	return tmpDir, stackFile.Name()
}

func cleanup(dir string) {
	os.RemoveAll(dir)
}

func TestNewDeployConfig(t *testing.T) {
	tmpDir, stackFile := setup(t)

	dc, err := NewDeployConfig(stackFile)
	assert.NoError(t, err)
	assert.IsType(t, dc, new(DeployConfig))

	cleanup(tmpDir)
}

func TestGetStackList(t *testing.T) {
	tmpDir, stackFile := setup(t)

	dc, _ := NewDeployConfig(stackFile)

	// Test name filer
	f := map[string]string{"name": "stack-a"}
	sc := dc.GetStackList(f)
	assert.Equal(t, 1, len(sc))

	f = map[string]string{"name": "stack-a,stack-b"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 2, len(sc))

	// Test tag filter
	f = map[string]string{"tag": "Name=stack-a"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 1, len(sc))

	f = map[string]string{"tag": "Name=stack-a,App=test"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 1, len(sc))

	f = map[string]string{"tag": "App=test"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 2, len(sc))

	// Test all filter
	sc = dc.GetStackList(nil)
	assert.Equal(t, 2, len(sc))

	// Test name and tag filter together
	f = map[string]string{"tag": "Name=stack-a", "name": "stack-a"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 1, len(sc))

	f = map[string]string{"tag": "Name=stack-a", "name": "stack-b"}
	sc = dc.GetStackList(f)
	assert.Equal(t, 0, len(sc))

	cleanup(tmpDir)
}
