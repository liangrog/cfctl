package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCipher(t *testing.T) {
	plainText := "ABCDEabcde12345"
	password := "test"

	v := new(Vault)

	encrypted, err := v.Encrypt([]byte(plainText), password)
	assert.NoError(t, err)
	assert.True(t, len(encrypted) > 0)

	decrypted, err := v.Decrypt(password, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plainText, string(decrypted))
}
