---
date: 2020-09-19T13:39:18+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Go 执行终端命令/外部命令"  # 文章标题
url:  "posts/go/libraries/tripartite/exec"  # 设置网页链接，默认使用文件名
tags: [ "go", "exec" ]  # 自定义标签
series: [ "Go 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
toc: false  # 是否自动生成目录
---

## 改变执行程序的环境

```go
cmd := exec.Command("programToExecute")
additionalEnv := "FOO=bar"
newEnv := append(os.Environ(), additionalEnv))
cmd.Env = newEnv
out, err := cmd.CombinedOutput()
if err != nil {
	log.Fatalf("cmd.Run() failed with %s\n", err)
}
fmt.Printf("%s", out)
```

```go

```

```go

```


## 忽略输出结果

```go
func main() {
	cmd := "python -m this"
	args := strings.Split(cmd, " ")
	_ = exec.Command(args[0], args[1:]...).Run()
}
```

其中 `Run()` 是对两个命令的封装：

```go
func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}
```

## 获取输出结果

```go
func main() {
	s := "python -m this"
	args := strings.Split(s, " ")
	c := exec.Command(args[0], args[1:]...)
	out, _ := c.CombinedOutput()
	fmt.Println(string(out))
}
```

## 分开处理错误输出和标准输出

```go
func main() {
	s := "python -m this.py"
	args := strings.Split(s, " ")
	c := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	_ = c.Run()
	fmt.Println("---stdout---")
	fmt.Println(stdout.String())
	fmt.Println("---stderr---")
	fmt.Println(stderr.String())
}
```

## 命令执行过程中获得输出

```go
func main() {
	s := "python -m this.py"
	args := strings.Split(s, " ")
	c := exec.Command(args[0], args[1:]...)
	stdoutIn, _ := c.StdoutPipe()
	stderrIn, _ := c.StderrPipe()
	var stdoutBuf, stderrBuf bytes.Buffer

	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	_ = c.Start()
	go func() { _, _ = io.Copy(stdout, stdoutIn) }()
	go func() { _, _ = io.Copy(stderr, stderrIn) }()
	_ = c.Wait()

	fmt.Println("---stdout---")
	fmt.Println(stdoutBuf.String())
	fmt.Println("---stderr---")
	fmt.Println(stderrBuf.String())
}
```

```
https://colobu.com/2017/06/19/advanced-command-execution-in-Go-with-os-exec
```


## 无法运行的情况

- 参数不可以带引号
- 必须分割命令参数

```go

```

```go

```

```go

```

```go

```

```go

```

```go

```
