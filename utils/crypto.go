package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
)

func evpBytesToKey(password, salt []byte) (key, iv []byte) {
	var dx = []byte{}
	hash := md5.New()
	for len(dx) < 48 { // AES-256 key(32) + IV(16)
		hash.Write(dx)
		hash.Write(password)
		hash.Write(salt)
		dx = hash.Sum(nil)
		hash.Reset()
	}
	key = dx[:32]
	iv = dx[32:48]
	return
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: data is empty")
	}
	padLen := int(data[length-1])
	if padLen > aes.BlockSize || padLen == 0 {
		return nil, errors.New("pkcs7: invalid padding")
	}
	return data[:length-padLen], nil
}

func DecryptData(encryptedBase64 string, passphrase string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}

	if len(cipherData) < 16 {
		return "", errors.New("cipher data too short")
	}

	if string(cipherData[:8]) != "Salted__" {
		return "", errors.New("missing OpenSSL salt header (Salted__)")
	}

	salt := cipherData[8:16]
	key, iv := evpBytesToKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := cipherData[16:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	unpadded, err := pkcs7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(unpadded, "\x00")), nil
}
