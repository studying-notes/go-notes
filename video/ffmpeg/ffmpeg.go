package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShowStream 显示流信息
func ShowStream(src string) (err error) {
	if !IsExist(src) {
		return IsNotExistError(src)
	}
	return Command(fmt.Sprintf("-i %s", src))
}

// ExtractVideoAudio 分离出无声视频和 MP3 音频
func ExtractVideoAudio(src string) (err error) {
	if err = ExtractVideo(src); err != nil {
		return err
	}
	if err = VideoConvAudio(src, src[:strings.LastIndex(src,
		".")]+"_audio.mp3"); err != nil {
		return err
	}
	return nil
}

// VideoConvAudio 视频转换为音频，以文件名后缀区分格式
func VideoConvAudio(src, dst string) (err error) {
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

// OverlayOriginAudio 替换原视频中的音频流
func OverlayOriginAudio(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", src, audio, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	c := fmt.Sprintf("-i %s -i %s -shortest -c:v copy -c:a aac "+
		"-strict experimental -map 0:v:0 -map 1:a:0 %s", src, audio, dst)
	return Command(c)
}

// MergeVideoAudio 混合不带音频的视频文件和音频文件
func MergeVideoAudio(video, audio, dst string) (err error) {
	if !IsExist(video) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", video, audio, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	// 音频或视频流较长，可以添加 -shortest 选项，以便一个文件结束后停止编码
	c := fmt.Sprintf("-i %s -i %s -shortest -c:v copy -c:a aac "+
		"-strict experimental %s", video, audio, dst)
	return Command(c)
}

// ReplaceOriginAudio 替换原视频中的音频流
// 此方法是先分离视频流再添加音频流
func ReplaceOriginAudio(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", src, audio, filepath.Dir(dst)))
	}
	// 分离视频流
	if err = ExtractVideo(src); err != nil {
		return err
	}
	video := src[:strings.LastIndex(src, ".")] + "_noaudio.mp4"

	// 清除中间文件
	defer func() {
		_ = os.Remove(video)
	}()

	return MergeVideoAudio(video, audio, dst)
}

// MixBackgroundMusic 混合视频原声和背景音乐，目前最佳方法
func MixBackgroundMusic(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", src, audio, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	if err = ExtractVideoAudio(src); err != nil {
		return err
	}

	video := src[:strings.LastIndex(src, ".")] + "_noaudio.mp4"
	mp3 := src[:strings.LastIndex(src, ".")] + "_audio.mp3"
	mix := src[:strings.LastIndex(src, ".")] + "_mix.mp3"
	if err = Mix2Audios(mp3, audio, mix); err != nil {
		return err
	}

	defer func() {
		_ = os.Remove(video)
		_ = os.Remove(mp3)
		_ = os.Remove(mix)
	}()

	return MergeVideoAudio(video, mix, dst)
}

// ExtractVideoImages 从视频提取图片
//func ExtractVideoImages(src string, frequency int) (err error) {
//
//}
