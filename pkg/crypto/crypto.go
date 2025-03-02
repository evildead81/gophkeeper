package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с помощью bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash проверяет соответствие пароля его хешу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateAESKey генерирует 32-байтный ключ для AES-256
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256 бит = 32 байта
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt шифрует данные с использованием AES-256
func Encrypt(plainText string) (string, error) {
	key := []byte("myverystrongpasswordo32bitlength") // В реальном проекте - хранить в конфиге или .env
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainData := []byte(plainText)
	cipherText := make([]byte, aes.BlockSize+len(plainData))
	iv := cipherText[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainData)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt расшифровывает данные с использованием AES-256
func Decrypt(cipherText string) (string, error) {
	key := []byte("myverystrongpasswordo32bitlength") // В реальном проекте - хранить в конфиге или .env
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("шифрованный текст слишком короткий")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return string(data), nil
}
