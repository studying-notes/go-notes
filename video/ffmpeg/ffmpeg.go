package ffmpeg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExtractVideo 分离视频流
func ExtractVideo(src string) (err error) {
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	dst := src[:strings.LastIndex(src, ".")] + "_noaudio.mp4"
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	c := fmt.Sprintf("-i %s -vcodec copy -an %s", src, dst)
	return Command(c)
}

// ExtractAudio 分离音频流
func ExtractAudio(src string) (err error) {
	if !IsExist(src) {
		return errors.New(fmt.Sprintf("%s isn't exist", src))
	}
	dst := src[:strings.LastIndex(src, ".")] + "_audio.m4a"
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	c := fmt.Sprintf("-i %s -acodec copy -vn %s", src, dst)
	return Command(c)
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

// OverlayOriginAudio 替换原视频中的音频流
func OverlayOriginAudio(src, audio, dst string) (err error) {
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

// MergeVideoAudio 混合不带音频的视频文件和音频文件
func MergeVideoAudio(video, audio, dst string) (err error) {
	if !IsExist(video) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s isn't exist", video, audio, filepath.Dir(dst)))
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
			"%s isn't exist", src, audio, filepath.Dir(dst)))
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
			"%s isn't exist", src, audio, filepath.Dir(dst)))
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

// AddBackgroundMusic 保留原声，添加背景音乐，有可能声音错位，不推荐
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
	c := fmt.Sprintf(mp, audio, src, dst) // 音频和视频输入流的命令顺序对视频合成有影响，顺序交换可能造成跳帧
	return Command(c)
}

// ConcatVideos 拼接多个 MP4 格式的文件
func ConcatVideos(dst string, videos ...string) (err error) {
	if !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s isn't exist", filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	ts := "-i %s -vcodec copy -acodec copy -vbsf h264_mp4toannexb %d.ts"
	var join []string
	// 不可以带引号，否则无法运行
	concat := "-i concat:%s -acodec copy -vcodec copy -absf aac_adtstoasc %s"
	for idx, video := range videos {
		if !IsExist(video) {
			return errors.New(fmt.Sprintf("%s isn't exist", video))
		}
		if ext := filepath.Ext(video); ext != ".mp4" {
			return errors.New("unsupported format: " + ext)
		}
		// 必须分割命令参数，否则无法运行
		args := strings.Split(fmt.Sprintf(ts, video, idx), " ")
		// 为了兼容性，将 mp4 格式都转换成 ts 格式
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

// Mix2Audios 将两个 MP3 文件混合合成一个，不是拼接
func Mix2Audios(x, y, dst string) (err error) {
	if !IsExist(x) || !IsExist(y) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s isn't exist", x, y, filepath.Dir(dst)))
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

// IsExist 判断文件/目录是否存在
func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}

// Command 命令执行入口
func Command(s string) (err error) {
	args := strings.Split(s, " ")
	c := exec.Command("ffmpeg", args...)

	// 命令执行过程中获得输出
	stdoutIn, _ := c.StdoutPipe()
	stderrIn, _ := c.StderrPipe()

	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	_ = c.Start()
	go func() { _, _ = io.Copy(stdout, stdoutIn) }()
	go func() { _, _ = io.Copy(stderr, stderrIn) }()

	if err = c.Wait(); err != nil {
		fmt.Println(c.String())
	}

	fmt.Println(stdoutBuf.String())
	fmt.Println(stderrBuf.String())

	return err
}
