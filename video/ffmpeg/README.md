# ffmpeg 学习笔记

+ [术语](term.md)
+ [常用命令](usage.md)

## Go 外部调用命令

### 合并多个 MP4 格式的文件

```go
// ConcatMP4 合并多个 MP4 格式的文件
func ConcatMP4(dst string, mp4s ...string) (err error) {
	ts := "-i %s -vcodec copy -acodec copy -vbsf h264_mp4toannexb %d.ts"
	var join []string
	// 不可以带引号，否则无法运行
	concat := "-i concat:%s -acodec copy -vcodec copy -absf aac_adtstoasc %s"
	for idx, mp4 := range mp4s {
		ext := filepath.Ext(mp4)
		if ext == "" {
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
	args := strings.Split(fmt.Sprintf(concat, strings.Join(join, "|"), dst), " ")
	// 合并全部 ts 文件为一个 mp4 文件
	return exec.Command("ffmpeg", args...).Run()
}
```
