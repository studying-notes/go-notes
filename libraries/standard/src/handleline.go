/*
	一行一行读取文件，然后对每一行进行处理
*/

package io

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func handle(line []byte) []byte {
	list := bytes.SplitN(line, []byte(" "), 2)
	return list[1]
}

func HandleLine(src, dst string, handle func([]byte) []byte) error {
	fr, err := os.Open(src)
	defer fr.Close()
	if err != nil {
		return err
	}

	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	bufReader := bufio.NewReader(fr)
	bufWriter := bufio.NewWriter(fw)

	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		}
		newLine := handle(line)
		_, _ = bufWriter.Write(newLine)
		_, _ = bufWriter.Write([]byte("\n"))
	}
	return bufWriter.Flush()
}
