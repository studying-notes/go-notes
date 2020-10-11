package io

import (
	"bytes"
	"github/fujiawei-dev/go-notes/encrypt"
	"testing"
)

// go test -run TestHandleLine
func TestHandleLine(t *testing.T) {
	src := "dev/encrypt.log"
	dst := "dev/decrypt.log"
	err := HandleLine(src, dst, func(r []byte) []byte {
		rs := bytes.Fields(r)
		rx, _ := encrypt.AesBase64Decrypt(string(rs[len(rs)-1]))
		return append(r[:len(r)-len(rs[len(rs)-1])], []byte(rx+"\n")...)
	})
	if err != nil {
		t.Error(err)
	}
}
