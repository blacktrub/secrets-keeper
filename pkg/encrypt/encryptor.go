package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

const defaultSecretKey = "supersecretkey11"


func getSecretKey() string {
    envSecret := os.Getenv("SECRET_KEY")
    if envSecret == "" {
        return defaultSecretKey
    }
    return envSecret
}

func Encrypt(msg string) (string, error) {
	c, err := aes.NewCipher([]byte(getSecretKey()))
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

	encMsg := gcm.Seal(nonce, nonce, []byte(msg), nil)
	return string(encMsg), nil
}

func Decrypt(msg string) (string, error) {
	c, err := aes.NewCipher([]byte(getSecretKey()))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(msg) < nonceSize {
		return "", err
	}

	nonce, encyptedMsg := []byte(msg[:nonceSize]), []byte(msg[nonceSize:])
	realMsg, err := gcm.Open(nil, nonce, encyptedMsg, nil)
	if err != nil {
		return "", err
	}
	return string(realMsg), nil
}
