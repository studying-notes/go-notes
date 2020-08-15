package ffmpeg

import (
	"fmt"
	"testing"
)

//  go test -run TestGetDuration
func TestGetDuration(t *testing.T) {
	d, err := GetDuration("mp4/1.mp4")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(d)
}
