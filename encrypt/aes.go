/*
	Advanced Encryption Standard
*/

package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var aesKey = []byte("EcSL4BgiCBqsy2Ct")

// RKCS7 Padding
func PKCS7Padding(cip []byte, blockSize int) []byte {
	padCount := blockSize - len(cip)%blockSize
	padBytes := bytes.Repeat([]byte{byte(padCount)}, padCount)
	return append(cip, padBytes...)
}

func PKC7UnPadding(cip []byte) ([]byte, error) {
	length := len(cip)
	if length == 0 {
		return nil, errors.New("[]byte is empty")
	} else {
		padCount := int(cip[length-1])
		return cip[:length-padCount], nil
	}
}

// AES 加密
func AesEncrypt(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	src = PKCS7Padding(src, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(src))
	blockMode.CryptBlocks(encrypted, src)
	return encrypted, nil
}

// AES 解密
func AesDecrypt(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(src))
	blockMode.CryptBlocks(decrypted, src)
	decrypted, err = PKC7UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// AES + Base64 加密
func AesBase64Encrypt(src string) (string, error) {
	encrypted, err := AesEncrypt([]byte(src), aesKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// AES + Base64 解密
func AesBase64Decrypt(src string) (string, error) {
	cip, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	decrypted, err := AesDecrypt(cip, aesKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
