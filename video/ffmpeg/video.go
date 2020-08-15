package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

// GetDuration 获取视频时长
func GetDuration(src string) (string, error) {
	if !IsExist(src) {
		return "00:00:00:00", IsNotExistError(src)
	}

	// 必须忽略错误，因为没有指定输出文件
	buf, _ := CombinedOutput(fmt.Sprintf("-i %s", src))

	reg, err := regexp.Compile(`Duration: (\d{2}:\d{2}:\d{2}.\d{2}), start`)
	if err != nil {
		return "00:00:00:00", err
	}

	// 第1个匹配到的是这个字符串本身，从第2个开始，才是想要的字符串
	return reg.FindStringSubmatch(string(buf))[1], nil
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
func CutVideo(src, dst string, start, length float32) (err error) {
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

// ApplyFadeSecond 按秒设置淡入淡出
func ApplyFadeSecond(src, dst string, in, ind, out, outd float32) (err error) {
	if err = IsAllExist(src, filepath.Dir(dst)); err != nil {
		return err
	}

	var vf string
	if in >= 0 && ind > 0 && out >= 0 && outd > 0 {
		vf = fmt.Sprintf("fade=in:st=%f:d=%f,fade=out:st=%f:d=%f", in, ind, out, outd)
	} else if in >= 0 && ind > 0 {
		vf = fmt.Sprintf("fade=in:st=%f:d=%f", in, ind)
	} else if out >= 0 && outd > 0 {
		vf = fmt.Sprintf("fade=out:st=%f:d=%f", out, outd)
	} else {
		return errors.New("invalid parameters")
	}

	return Command(fmt.Sprintf("-i %s -vf %s %s -y", src, vf, dst))
}

func ConcatVideosApplyFade(dst string, videos ...string) (err error) {
	if !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s/ isn't exist", filepath.Dir(dst)))
	}
	var vs []string

	for idx, video := range videos {
		if !IsExist(video) {
			return errors.New(fmt.Sprintf("%s isn't exist", video))
		}
		if ext := filepath.Ext(video); ext != ".mp4" {
			return errors.New("unsupported format: " + ext)
		}

		d, err := GetDuration(video)
		if err != nil {
			return err
		}
		length := ConvertString(d)
		fmt.Println(length)
		v := fmt.Sprintf("%d_vs.mp4", idx)
		if err = ApplyFadeSecond(video, v, 0, 0.2, length-0.2, 0.2); err != nil {
			return err
		}
		vs = append(vs, v)
	}
	// 清除中间文件
	defer func() {
		for _, v := range vs {
			_ = os.Remove(v)
		}
	}()
	return ConcatVideos(dst, vs...)
}
