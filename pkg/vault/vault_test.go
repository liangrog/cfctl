package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCipher(t *testing.T) {
	plainText := "ABCDEabcde12345"
	password := "test"

	v := new(Vault)

	err := v.Encrypt([]byte(plainText), password)
	assert.NoError(t, err)

	encrypted := v.Encode()

	decrypted, err := v.Decrypt(password, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plainText, string(decrypted))
}
