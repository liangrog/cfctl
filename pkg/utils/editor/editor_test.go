/*
File Source: kubernetes/staging/src/k8s.io/kubectl/pkg/cmd/util/editor/editor_test.go
*/

package editor

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	if e, a := []string{"/bin/bash", "-c \"test\""}, (Editor{Args: []string{"/bin/bash", "-c"}, Shell: true}).args("test"); !reflect.DeepEqual(e, a) {
		t.Errorf("unexpected args: %v", a)
	}
	if e, a := []string{"/bin/bash", "-c", "test"}, (Editor{Args: []string{"/bin/bash", "-c"}, Shell: false}).args("test"); !reflect.DeepEqual(e, a) {
		t.Errorf("unexpected args: %v", a)
	}
	if e, a := []string{"/bin/bash", "-i -c \"test\""}, (Editor{Args: []string{"/bin/bash", "-i -c"}, Shell: true}).args("test"); !reflect.DeepEqual(e, a) {
		t.Errorf("unexpected args: %v", a)
	}
	if e, a := []string{"/test", "test"}, (Editor{Args: []string{"/test"}}).args("test"); !reflect.DeepEqual(e, a) {
		t.Errorf("unexpected args: %v", a)
	}
}

func TestEditor(t *testing.T) {
	edit := Editor{Args: []string{"cat"}}
	testStr := "test something\n"
	contents, path, err := edit.LaunchTempFile("", "someprefix", bytes.NewBufferString(testStr))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("no temp file: %s", path)
	}
	defer os.Remove(path)
	if disk, err := ioutil.ReadFile(path); err != nil || !bytes.Equal(contents, disk) {
		t.Errorf("unexpected file on disk: %v %s", err, string(disk))
	}
	if !bytes.Equal(contents, []byte(testStr)) {
		t.Errorf("unexpected contents: %s", string(contents))
	}
	if !strings.Contains(path, "someprefix") {
		t.Errorf("path not expected: %s", path)
	}
}
