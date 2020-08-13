package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExtractVideo 分离视频流
func ExtractVideo(src string) (err error) {
	if !IsExist(src) {
		return IsNotExistError(src)
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
		return IsNotExistError(src)
	}
	dst := src[:strings.LastIndex(src, ".")] + "_audio.m4a"
	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	c := fmt.Sprintf("-i %s -acodec copy -vn %s", src, dst)
	return Command(c)
}

// ConcatVideos 拼接多个 MP4 格式的文件
func ConcatVideos(dst string, videos ...string) (err error) {
	if !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s/ isn't exist", filepath.Dir(dst)))
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
		// 为了兼容性，将 mp4 格式都转换成 ts 格式
		if err = Command(fmt.Sprintf(ts, video, idx)); err != nil {
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

// CutVideo 截取视频片段，开始位置和截取长度
func CutVideo(src, dst string, start, length int) (err error) {
	if !IsExist(src) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s/ isn't exist", src, filepath.Dir(dst)))
	}
	ss := "-ss " + ConvertSecond(start)
	if length > 0 { // length <= 0 时只去除片头
		ss += " -t " + ConvertSecond(length)
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	return Command(fmt.Sprintf("-i %s %s -vcodec copy -acodec copy %s", src, ss, dst))
}

// TransposeVideo 旋转视频
func TransposeVideo(src, dst string, angle int) (err error) {
	if err = IsAllExist(src, filepath.Dir(dst)); err != nil {
		return err
	}

	var vf string
	switch angle {
	case 90:
		vf = "transpose=1" // 顺时针水平旋转 90 度
	case -90:
		vf = "transpose=2" // 逆时针水平旋转 90 度
	case 180:
		vf = "transpose=1,transpose=1"
	default:
		return errors.New("only support -90, 90 and 180")
	}

	// 移除可能存在的目标文件
	_ = os.Remove(dst)

	return Command(fmt.Sprintf("-i %s -vf %s %s", src, vf, dst))
}
