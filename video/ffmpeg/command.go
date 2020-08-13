package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Command 命令执行入口
func Command(s string) (err error) {

	// 必须分割命令参数，否则无法运行
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
