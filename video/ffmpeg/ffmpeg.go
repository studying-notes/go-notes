package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	err = ConcatVideos("mp4/merge.mp4", "mp4/1.mp4",
		"mp4/2.mp4", "mp4/3.mp4", "mp4/4.mp4", "mp4/5.mp4")
	if err != nil {
		panic(err)
	}
	err = AddBackgroundMusic("mp4/merge.mp4",
		"mp4/1.mp3", "mp4/bgm.mp4")
	if err != nil {
		panic(err)
	}
}

// ExtractVideo 分离视频流
func ExtractVideo(src string) (err error) {
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	c := fmt.Sprintf("-i %s -vcodec copy -an %s_noaudio.mp4",
		src, src[:strings.LastIndex(src, ".")])
	return Command(c)
}

// ExtractAudio 分离音频流
func ExtractAudio(src string) (err error) {
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	c := fmt.Sprintf("-i %s -acodec copy -vn %s_audio.m4a",
		src, src[:strings.LastIndex(src, ".")])
	return Command(c)
}

// AudioFormatConv 音频格式转换，以文件名后缀区分格式
func AudioFormatConv(src, dst string) (err error) {
	if src == dst {
		return errors.New("source file and target file have the same name")
	}
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	c := fmt.Sprintf("-i %s %s", src, dst)
	return Command(c)
}

// VideoConvAudio 视频转换为音频，以文件名后缀区分格式
func VideoConvAudio(src, dst string) (err error) {
	if src == dst {
		return errors.New("source file and target file have the same name")
	}
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	c := fmt.Sprintf("-i %s %s", src, dst)
	return Command(c)
}

// ReplaceOriginAudio 替换原视频中的音频流
func ReplaceOriginAudio(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s isn't exist", src, audio, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	c := fmt.Sprintf("-i %s -i %s -shortest -c:v copy -c:a aac "+
		"-strict experimental -map 0:v:0 -map 1:a:0 %s", src, audio, dst)
	return Command(c)
}

//func ReplaceOriginAudio(src, audio, dst string) (err error) {
//	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
//		return errors.New(fmt.Sprintf("%s or %s or "+
//			"%s isn't exist", src, audio, filepath.Dir(dst)))
//	}
//	// 移除可能存在的目标文件
//	_ = os.Remove(dst)
//
//	// 分离视频流
//	if err = ExtractVideo(src); err != nil {
//		return err
//	}
//	noaudio := src[:strings.LastIndex(src, ".")] + "_noaudio.mp4"
//	// 清除中间文件
//	defer func() {
//		_ = os.Remove(noaudio)
//	}()
//
//	// 音频或视频流较长，可以添加 -shortest 选项，以便一个文件结束后停止编码
//	c := fmt.Sprintf("-i %s -i %s -shortest -c:v copy -c:a aac "+
//		"-strict experimental %s", noaudio, audio, dst)
//	return Command(c)
//}

// AddBackgroundMusic 保留原声，添加背景音乐
func AddBackgroundMusic(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s isn't exist", src, audio, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	// inputs：输入流数量；duration：决定流的结束
	// dropout_transition：输入流结束时，容量重整时间
	// longest 最长输入时间，shortest 最短，first 第一个输入持续的时间
	mp := "-i %s -i %s -threads 16 -filter_complex amix=inputs=2:duration=longest:dropout_transition=0 %s -y"
	c := fmt.Sprintf(mp, audio, src, dst) // 音频/视频输入流的命令顺序对视频合成有影响，顺序交换可能造成跳帧
	return Command(c)
}

// ConcatVideos 合并多个 MP4 格式的文件
func ConcatVideos(dst string, mp4s ...string) (err error) {
	if !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s isn't exist", filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	ts := "-i %s -vcodec copy -acodec copy -vbsf h264_mp4toannexb %d.ts"
	var join []string
	// 不可以带引号，否则无法运行
	concat := "-i concat:%s -acodec copy -vcodec copy -absf aac_adtstoasc %s"
	for idx, mp4 := range mp4s {
		if !IsExist(mp4) {
			return errors.New(fmt.Sprintf("%s isn't exist", mp4))
		}
		if ext := filepath.Ext(mp4); ext != ".mp4" {
			return errors.New("unsupported format: " + ext)
		}
		// 必须分割命令参数，否则无法运行
		args := strings.Split(fmt.Sprintf(ts, mp4, idx), " ")
		// 将 mp4 格式都转换成 ts 格式
		if err = exec.Command("ffmpeg", args...).Run(); err != nil {
			return err
		}
		join = append(join, fmt.Sprintf("%d.ts", idx))
	}
	// 清除中间文件
	defer func() {
		for _, v := range join {
			_ = os.Remove(v)
		}
	}()
	if ext := filepath.Ext(dst); ext == "" {
		dst += ".mp4"
	} else if ext != ".mp4" {
		return errors.New("unsupported format: " + ext)
	}
	c := fmt.Sprintf(concat, strings.Join(join, "|"), dst)
	// 将 ts 文件合并为一个 mp4 文件
	return Command(c)
}

// IsExist 判断文件/目录是否存在
func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}

// Command 命令执行入口
func Command(c string) error {
	args := strings.Split(c, " ")
	return exec.Command("ffmpeg", args...).Run()
}
