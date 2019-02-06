package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Ansible vault version
	header = "$ANSIBLE_VAULT;1.1;AES256"

	// Spec
	cipherKeyLength = 32
	HMACKeyLength   = 32
	saltLength      = 32
	ivLength        = 16
	iteration       = 10000

	// Format
	charPerLine = 80
)

// Key used to cipher
type cipherKey struct {
	key     []byte
	hmacKey []byte
	iv      []byte
}

// Generate cipher key
func keyGen(password string, salt []byte) *cipherKey {
	k := pbkdf2.Key(
		[]byte(password),
		salt,
		iteration,
		(cipherKeyLength + HMACKeyLength + ivLength),
		sha256.New,
	)

	return &cipherKey{
		key:     k[:cipherKeyLength],
		hmacKey: k[cipherKeyLength:(cipherKeyLength + HMACKeyLength)],
		iv:      k[(cipherKeyLength + HMACKeyLength):(cipherKeyLength + HMACKeyLength + ivLength)],
	}
}

// Generate salt
func saltGen(n int) ([]byte, error) {
	s := make([]byte, n)
	_, err := rand.Read(s)

	return s, err
}

// Encrypt or decrypt the given data
func cipherData(cipherType string, data []byte, key *cipherKey) ([]byte, error) {
	var output []byte

	block, err := aes.NewCipher(key.key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, key.iv)

	switch cipherType {
	case "encrypt":
		data = aesBlockPad(data)
		output = make([]byte, len(data))
		stream.XORKeyStream(output, data)
	case "decrypt":
		decryptedData := make([]byte, len(data))
		stream.XORKeyStream(decryptedData, data)
		output, err = aesBlockUnpad(decryptedData)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Missing instruction on how to cipher")
	}

	return output, nil
}

// Pad data to AES block
func aesBlockPad(data []byte) []byte {
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	return append(data, (bytes.Repeat([]byte{byte(padLen)}, padLen))...)
}

// Unpad data for AES block
func aesBlockUnpad(data []byte) ([]byte, error) {
	length := len(data)
	unpad := int(data[length-1])

	if unpad > length {
		return nil, errors.New("Unpad error")
	}

	return data[:(length - unpad)], nil
}

// Format the text to given length
func wrap(data []byte, length int) string {
	var output []byte

	for i := 0; i < len(data); i++ {
		// Append prune line wrap char
		if i > 0 && i%length == 0 {
			output = append(output, '\n')
		}
		output = append(output, data[i])
	}

	return string(output)
}

// Validate hmac checksum
func validateCheckSum(checkSum, data, hmacKey []byte) bool {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(data)
	return hmac.Equal(mac.Sum(nil), checkSum)
}

// Data to be encrypted
type Vault struct {
	data     []byte
	checkSum []byte
	salt     []byte
}

// Encypt given data for the password
func (v *Vault) Encrypt(data []byte, password string) ([]byte, error) {
	var err error

	// Empty password is not allowed
	if len(password) <= 0 {
		return nil, errors.New("Empty password")
	}

	lines := strings.SplitN(string(data), "\n", 2)
	if strings.TrimSpace(lines[0]) == header {
		return nil, errors.New("Given data has already been encrypted according to header")
	}

	// Get salt
	v.salt, err = saltGen(saltLength)
	if err != nil {
		return nil, err
	}

	key := keyGen(password, v.salt)

	v.data, err = cipherData("encrypt", data, key)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, key.hmacKey)
	mac.Write(v.data)
	v.checkSum = mac.Sum(nil)

	return v.Encode(), nil
}

func (v *Vault) Encode() []byte {
	content := []byte(
		strings.Join(
			[]string{
				hex.EncodeToString(v.salt),
				hex.EncodeToString(v.checkSum),
				hex.EncodeToString(v.data),
			},
			"\n",
		))

	return []byte(strings.Join(
		[]string{
			header,
			wrap([]byte(hex.EncodeToString(content)), charPerLine),
		}, "\n"))
}

func (v *Vault) Decode(str string) error {
	lines := strings.SplitN(str, "\n", 2)

	if strings.TrimSpace(lines[0]) != header {
		return errors.New("Invalid vault file format")
	}

	// Concat all lines
	content := strings.TrimSpace(lines[1])
	content = strings.Replace(content, "\r", "", -1)
	content = strings.Replace(content, "\n", "", -1)

	// Decode the first layer
	decodedStr, err := hex.DecodeString(content)
	if err != nil {
		return err
	}

	lines = strings.Split(string(decodedStr), "\n")
	if len(lines) != 3 {
		return errors.New("Invalid encoded data")
	}

	if v.salt, err = hex.DecodeString(lines[0]); err != nil {
		return err
	}

	if v.checkSum, err = hex.DecodeString(lines[1]); err != nil {
		return err
	}

	if v.data, err = hex.DecodeString(lines[2]); err != nil {
		return err
	}

	return nil
}

func (v *Vault) Decrypt(password string, data []byte) ([]byte, error) {
	// Empty password is not allowed
	if len(password) <= 0 {
		return nil, errors.New("Empty password")
	}

	if err := v.Decode(string(data)); err != nil {
		return nil, err
	}

	key := keyGen(password, v.salt)

	return cipherData("decrypt", v.data, key)
}
