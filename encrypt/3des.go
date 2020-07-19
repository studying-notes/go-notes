/*
	Data Encryption Standard
*/

package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
)

// 3DES 的秘钥长度必须为24位
var desKey = []byte("123456781234567812345678")

/*
	The difference between the PKCS#5 and PKCS#7
	padding mechanisms is the block size; PKCS#5
	padding is defined for 8-byte block sizes, PKCS#7
	padding would work for any block size from 1 to 255 bytes.
*/

// 此处实现与 RKCS7 是一样的
// RKCS5 Padding
func PKCS5Padding(cip []byte, blockSize int) []byte {
	padCount := blockSize - len(cip)%blockSize
	padBytes := bytes.Repeat([]byte{byte(padCount)}, padCount)
	return append(cip, padBytes...)
}

func PKC5UnPadding(cip []byte) ([]byte, error) {
	length := len(cip)
	if length == 0 {
		return nil, errors.New("[]byte is empty")
	} else {
		padCount := int(cip[length-1])
		return cip[:length-padCount], nil
	}
}

// 3DES 加密
func TripleDesEncrypt(src []byte, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	src = PKCS5Padding(src, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(src))
	blockMode.CryptBlocks(encrypted, src)
	return encrypted, nil
}

// 3DES 解密
func TipleDesDecrypt(src []byte, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(src))
	blockMode.CryptBlocks(decrypted, src)
	decrypted, err = PKC5UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// DES + Base64 加密
func DesBase64Encrypt(src string) (string, error) {
	encrypted, err := TripleDesEncrypt([]byte(src), desKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DES + Base64 解密
func DesBase64Decrypt(src string) (string, error) {
	cip, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	decrypted, err := TipleDesDecrypt(cip, desKey)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
