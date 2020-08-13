package ffmpeg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Deprecated: AddBackgroundMusic should implement MixBackgroundMusic instead (or
// additionally). 保留原声，添加背景音乐，有可能声音错位，不推荐
func AddBackgroundMusic(src, audio, dst string) (err error) {
	if !IsExist(src) || !IsExist(audio) || !IsExist(filepath.Dir(dst)) {
		return errors.New(fmt.Sprintf("%s or %s or "+
			"%s/ isn't exist", src, audio, filepath.Dir(dst)))
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
