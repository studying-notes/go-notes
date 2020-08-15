package ffmpeg

import (
	"fmt"
	"testing"
)

// go test -run TestTruncateSecond
func TestTruncateSecond(t *testing.T) {
	fmt.Println(TruncateSecond("00:00:00.23"))
}

// go test -run TestConvertString
func TestConvertString(t *testing.T) {
	fmt.Println(ConvertString("00:02:03.23"))
}

// go test -run TestConvertSecond
func TestConvertSecond(t *testing.T) {
	fmt.Println(ConvertSecond(60.3))
}
