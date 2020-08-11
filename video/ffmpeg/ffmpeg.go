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
	err = ConcatMP4("mp4/merge.mp4", "mp4/1.mp4",
		"mp4/2.mp4", "mp4/3.mp4", "mp4/4.mp4", "mp4/5.mp4")
	if err != nil {
		panic(err)
	}
	err = AddBackgroundMusic("mp4/merge.mp4",
		"mp4/1.mp3", "mp4/background.mp4")
	if err != nil {
		panic(err)
	}
}

// AddBackgroundMusic 保留原声，添加背景音乐
func AddBackgroundMusic(src, mp3, dst string) (err error) {
	if !IsExist(src) || !IsExist(mp3) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s isn't exist", src, mp3, filepath.Dir(dst)))
	}
	// 移除可能存在的目标文件
	_ = os.Remove(dst)
	if ext := filepath.Ext(mp3); ext == "" {
		mp3 += ".mp3"
	} else if ext != ".mp3" {
		return errors.New("unsupported format: " + ext)
	}
	mp := "-i %s -i %s -threads 2 -filter_complex amix=inputs=2:duration=first:dropout_transition=0 %s -y"
	args := strings.Split(fmt.Sprintf(mp, src, mp3, dst), " ")
	return exec.Command("ffmpeg", args...).Run()
}

// ConcatMP4 合并多个 MP4 格式的文件
func ConcatMP4(dst string, mp4s ...string) (err error) {
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
		if ext := filepath.Ext(mp4); ext == "" {
			mp4 += ".mp4"
		} else if ext != ".mp4" {
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
	args := strings.Split(fmt.Sprintf(concat, strings.Join(join, "|"), dst), " ")
	// 合并全部 ts 文件为一个 mp4 文件
	return exec.Command("ffmpeg", args...).Run()
}

func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
