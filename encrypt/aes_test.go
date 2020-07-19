package encrypt

import (
	"fmt"
	"testing"
)

func TestAesBase64Encrypt(t *testing.T) {
	src := "雪落寂无声，少年亦萧瑟。"
	encrypted, err := AesBase64Encrypt(src)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(encrypted)
}

func TestAesBase64Decrypt(t *testing.T) {
	encrypted := "Hd+tGMQ3hqe2/m6zFtjroEJp3W2ECQ/q5P4HDSXMfzk6J37y8AdPb5mjTI3WQWtg"
	decrypted, err := AesBase64Decrypt(encrypted)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(decrypted)
}
