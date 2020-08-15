package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// AudioFormatConv 音频格式转换，以文件名后缀区分格式
func AudioFormatConv(src, dst string) (err error) {
	if src == dst {
		return errors.New("source file and target file have the same name")
	}
	if !IsExist(src) {
		return IsNotExistError(src)
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	c := fmt.Sprintf("-i %s %s", src, dst)
	return Command(c)
}

// Mix2Audios 将两个 MP3 文件混合合成一个，不是拼接
func Mix2Audios(x, y, dst string) (err error) {
	if !IsExist(x) || !IsExist(y) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", x, y, filepath.Dir(dst)))
	}
	if filepath.Ext(x) != ".mp3" || filepath.Ext(y) != ".mp3" || filepath.Ext(dst) != ".mp3" {
		return errors.New("only support MP3 format")
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	c := fmt.Sprintf("-i %s -i %s -filter_complex amerge -ac "+
		"2 -c:a libmp3lame -q:a 4 %s", x, y, dst)
	return Command(c)
}

// CutAudio 截取音频片段，开始位置和截取长度
func CutAudio(src, dst string, start, length float32) (err error) {
	if !IsExist(src) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s/ isn't exist", src, filepath.Dir(dst)))
	}
	if length <= 0 {
		return errors.New("length <=0 isn't supported")
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	return Command(fmt.Sprintf("-i %s -ss %s -t %s -acodec copy %s",
		src, ConvertSecond(start), ConvertSecond(length), dst))
}

// AdjustVolumeMultiple 倍数调整音频的音量
func AdjustVolumeMultiple(src, dst string, vol float32) (err error) {
	if err = IsAllExist(src, filepath.Dir(dst)); err != nil {
		return err
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	return Command(fmt.Sprintf("-i %s -filter:a volume=%f %s", src, vol, dst))
}

// AdjustVolumedB 加减 dB 调整音频的音量
func AdjustVolumedB(src, dst string, d int) (err error) {
	if err = IsAllExist(src, filepath.Dir(dst)); err != nil {
		return err
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	return Command(fmt.Sprintf("-i %s -filter:a volume=%ddB %s", src, d, dst))
}
