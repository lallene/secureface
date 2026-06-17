package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

func EncryptFaceVector(plainText string) (string, error) {
	key := []byte(os.Getenv("FACE_ENCRYPTION_KEY"))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptFaceVector(encryptedText string) (string, error) {
	key := []byte(os.Getenv("FACE_ENCRYPTION_KEY"))

	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("données chiffrées invalides")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func SaveEncryptedFaceVector(userID uint, plainText string) (string, error) {
	encryptedText, err := EncryptFaceVector(plainText)
	if err != nil {
		return "", err
	}

	dir := "storage/embeddings"
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/user_%d_face.enc", dir, userID)

	if err := os.WriteFile(filePath, []byte(encryptedText), 0600); err != nil {
		return "", err
	}

	return filePath, nil
}

func LoadEncryptedFaceVector(filePath string) (string, error) {
	encryptedText, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return DecryptFaceVector(string(encryptedText))
}
