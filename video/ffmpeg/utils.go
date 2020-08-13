package ffmpeg

import (
	"errors"
	"fmt"
	"os"
)

// ConvertSecond 将秒转换成时分秒 10:20:30 形式
func ConvertSecond(second int) string {
	return fmt.Sprintf("%02d:%02d:%02d",
		second/3600, second%3600/60, second%3600%60)
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
