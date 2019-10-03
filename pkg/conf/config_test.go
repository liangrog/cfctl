package conf

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/liangrog/vault"
	"github.com/stretchr/testify/assert"
)

func TestLoadVaules(t *testing.T) {
	// Setup files
	unixTime := time.Now().Unix()
	tmpDir, err := ioutil.TempDir("", "test"+strconv.FormatInt(unixTime, 10))
	assert.NoError(t, err)

	// plain text file
	pf, err := ioutil.TempFile(tmpDir, "plaintext")
	assert.NoError(t, err)
	_, err = pf.Write([]byte("plainVaule: 1234"))
	assert.NoError(t, err)

	// secret text file
	sf, err := ioutil.TempFile(tmpDir, "secret")
	assert.NoError(t, err)

	ss := "secretVaule: ABCD"
	pass := []string{"password"}
	ssb, err := vault.Encrypt([]byte(ss), pass[0])
	assert.NoError(t, err)
	_, err = sf.Write(ssb)
	assert.NoError(t, err)

	result, err := LoadValues(tmpDir, pass)
	assert.NoError(t, err)
	v, ok := result["plainVaule"]
	assert.True(t, ok)
	assert.Equal(t, "1234", v)

	v, ok = result["secretVaule"]
	assert.True(t, ok)
	assert.Equal(t, "ABCD", v)

	// test override
	sf, err = ioutil.TempFile(tmpDir, "tecret")
	assert.NoError(t, err)

	ss = "secretVaule: EFGH"
	ssb, err = vault.Encrypt([]byte(ss), pass[0])
	assert.NoError(t, err)
	_, err = sf.Write(ssb)
	assert.NoError(t, err)

	result, err = LoadValues(tmpDir, pass)
	assert.NoError(t, err)
	v, ok = result["secretVaule"]
	assert.True(t, ok)
	assert.Equal(t, "EFGH", v)

	// Clean up
	os.RemoveAll(tmpDir)
}
