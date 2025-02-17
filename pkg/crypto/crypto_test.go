package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionDecryption(t *testing.T) {
	original := "supersecretpassword"

	encrypted, err := Encrypt(original)
	assert.NoError(t, err)
	assert.NotEqual(t, original, encrypted)

	decrypted, err := Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, original, decrypted)
}
