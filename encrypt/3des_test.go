package encrypt

import (
	"fmt"
	"testing"
)

func TestDesBase64Encrypt(t *testing.T) {
	src := "雪落寂无声，少年亦萧瑟。"
	encrypted, err := DesBase64Encrypt(src)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(encrypted)
}

func TestDesBase64Decrypt(t *testing.T) {
	encrypted := "XzyPLkKQ/BH5re4dW5WuMQHkz23pbo75C05D+DWGA/sHZrCgIZs7LQ=="
	decrypted, err := DesBase64Decrypt(encrypted)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(decrypted)
}
