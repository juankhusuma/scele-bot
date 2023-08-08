package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func Encrypt(key string, val string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	seal := gcm.Seal(nonce, nonce, []byte(val), nil)
	return string(seal), nil
}

func Decrypt(key string, ci string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ci) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ci := ci[:nonceSize], ci[nonceSize:]
	val, err := gcm.Open(nil, []byte(nonce), []byte(ci), nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
