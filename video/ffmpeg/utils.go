package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ConvertSecond 将秒转换成时分秒
func ConvertSecond(second float32) string {
	hour := int(second) / 3600
	minute := int(second) % 3600 / 60
	second = second - float32(hour*3600+minute*60)
	return fmt.Sprintf("%02d:%02d:%.2f", hour, minute, second)
}

// ConvertString 将时分秒转换成秒
func ConvertString(s string) float32 {
	ts := strings.Split(s, ":")
	hour, _ := strconv.ParseFloat(ts[0], 32)
	minute, _ := strconv.ParseFloat(ts[1], 32)
	second, _ := strconv.ParseFloat(ts[2], 32)
	return float32(hour*3600 + minute*60 + second)
}

// TruncateSecond 截断到秒，丢弃小数
func TruncateSecond(s string) string {
	return s[:len(s)-3]
}

// IsExist 判断文件/目录是否存在
func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}

// IsNotExistError 返回文件不存在的错误信息
func IsNotExistError(src string) (err error) {
	return errors.New(fmt.Sprintf("%s isn't exist", src))
}

// IsAllExist 判断文件是否都存在
func IsAllExist(ps ...string) (err error) {
	for _, p := range ps {
		if !IsExist(p) {
			return IsNotExistError(p)
		}
	}
	return nil
}
