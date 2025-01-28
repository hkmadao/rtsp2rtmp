package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

func EncryptAES(key []byte, plainBytes []byte) (bytes []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// The IV needs to be unique, but does not have to be secret.
	// You can calculate it once with the key and keep it with the ciphertext.
	iv := key[:block.BlockSize()]

	stream := cipher.NewOFB(block, iv)

	bytes = make([]byte, len(plainBytes))
	stream.XORKeyStream(bytes, plainBytes)

	return
}

func DecryptAES(key []byte, encrypted []byte) (bytes []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	iv := key[:block.BlockSize()]

	stream := cipher.NewOFB(block, iv)

	bytes = make([]byte, len(encrypted))
	stream.XORKeyStream(bytes, encrypted)

	return
}
