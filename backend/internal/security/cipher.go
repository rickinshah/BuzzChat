package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"io"
)

func EncryptionWithHex(pt string) (string, error) {
	encryptionKey := getEncryptionKey()
	cipherKey := []byte(encryptionKey)

	if len(cipherKey) != 16 && len(cipherKey) != 24 && len(cipherKey) != 32 {
		return "", errors.New("invalid encryption key")
	}

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(pt), nil)

	final := append(nonce, cipherText...)

	return encode(final), nil

}

func DecryptionWithHex(ct string) (string, error) {
	encryptionKey := getEncryptionKey()
	cipherKey := []byte(encryptionKey)

	if len(cipherKey) != 16 && len(cipherKey) != 24 && len(cipherKey) != 32 {
		return "", errors.New("invalid encryption key")
	}

	data, err := decode(ct)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func encode(cipherText []byte) string {
	return base64.StdEncoding.EncodeToString(cipherText)
}

func decode(cipherText string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cipherText)
}

func getEncryptionKey() string {
	return flag.Lookup("encryption-key").Value.String()
}
