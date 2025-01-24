package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func EncryptAES(key []byte, plaintext string) (string, error) {
	plainBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but does not have to be secret.
	// You can calculate it once with the key and keep it with the ciphertext.
	iv := key[:block.BlockSize()]

	stream := cipher.NewOFB(block, iv)

	ciphertext := make([]byte, len(plainBytes))
	stream.XORKeyStream(ciphertext, plainBytes)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptAES(key []byte, ciphertext string) (string, error) {
	encrypted, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := key[:block.BlockSize()]

	stream := cipher.NewOFB(block, iv)

	plaintext := make([]byte, len(encrypted))
	stream.XORKeyStream(plaintext, encrypted)

	return string(plaintext), nil
}
