---
date: 2022-05-02T21:26:42+08:00
author: "Rustle Karl"  # 作者

title: "模糊测试 Fuzzing"  # 文章标题
url:  "posts/go/quickstart/feature/testing/fuzz"  # 设置网页永久链接
tags: [ "Go", "fuzz" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

> https://go.dev/doc/fuzz/

Fuzzing 是通过持续给一个程序不同的输入来自动化测试，并通过分析代码覆盖率来智能的寻找失败的 case。这种方法可以尽可能的寻找到一些边缘 case。

## 示例

### 实现一个函数

实现一个函数来对字符串做反转：

```go
package main

import "fmt"

func Reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	input := "The quick brown fox jumped over the lazy dog"
	rev := Reverse(input)
	doubleRev := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
```

### 增加单元测试

```go
package main

import (
	"testing"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, tc := range testcases {
		rev := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}
```

测试：

```shell
go test
```

### 增加模糊测试

单元测试有局限性，每个测试输入必须由开发者指定加到单元测试的测试用例里。

fuzzing 的优点之一是可以基于开发者代码里指定的测试输入作为基础数据，进一步自动生成新的随机测试数据，用来发现指定测试输入没有覆盖到的边界情况。

### 编写模糊测试

```go
package main

import (
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
```

Fuzzing 也有一定的局限性。

在单元测试里，因为测试输入是固定的，你可以知道调用 Reverse 函数后每个输入字符串得到的反转字符串应该是什么，然后在单元测试的代码里判断 Reverse 的执行结果是否和预期相符。例如，对于测试用例 Reverse("Hello, world")，单元测试预期的结果是 "dlrow ,olleH"。

但是使用 fuzzing 时，我们没办法预期输出结果是什么，因为测试的输入除了我们代码里指定的用例之外，还有 fuzzing 随机生成的。对于随机生成的测试输入，我们当然没办法提前知道输出结果是什么。

虽然如此，Reverse 函数有几个特性我们还是可以在模糊测试里做验证。

- 对一个字符串做2次反转，得到的结果和源字符串相同
- 反转后的字符串也仍然是一个有效的 UTF-8 编码的字符串

### 运行模糊测试

`go test` 方式只会使用种子语料库，而不会生成随机测试数据。

基于种子语料库生成随机测试数据用于模糊测试，需要给 go test 命令增加 -fuzz 参数。

```shell
go test -fuzz=Fuzz
```

上面的 fuzzing 测试结果是 FAIL，引起 FAIL 的输入数据被写到了一个语料库文件里。下次运行 go test 命令的时候，即使没有 -fuzz 参数，这个语料库文件里的测试数据也会被用到。

可以用文本编辑器打开 testdata/fuzz/FuzzReverse 目录下的文件，看看引起 Fuzzing 测试失败的测试数据长什么样。

```
go test fuzz v1
string("泃")
```

语料库文件里的第 1 行标识的是编码版本，虽然目前只有 v1 这 1 个版本，但是 Fuzzing 设计者考虑到未来可能引入新的编码版本，于是加了编码版本的概念。

从第 2 行开始，每一行数据对应的是语料库的每条测试数据 (corpus entry) 的其中一个参数，按照参数先后顺序排列。

本文的 FuzzReverse 里的 fuzz target 函数 func(t * testing.T, orig string) 只有 orig 这 1 个参数作为真正的测试输入，也就是每条测试数据其实就 1 个输入，因此在上面示例的 testdata/fuzz/FuzzReverse 目录下的文件里只有 string("泃") 这一行。

如果每条测试数据有 N 个参数，那 fuzzing 找出的导致 fuzz test 失败的每条测试数据在 testdata 目录下的文件里会有 N 行，第 i 行对应第 i 个参数。

### 问题原因

我们实现的 Reverse 函数是按照字节 (byte) 为维度进行字符串反转，这就是问题所在。比如中文里的字符 泃其实是由 3 个字节组成的，如果按照字节反转，反转后得到的就是一个无效的字符串了。

### 修复后的程序

```go
import (
    "errors"
    "fmt"
    "unicode/utf8"
)

func Reverse(s string) (string, error) {
    if !utf8.ValidString(s) {
        return s, errors.New("input is not valid UTF-8")
    }
    r := []rune(s)
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r), nil
}

func main() {
    input := "The quick brown fox jumped over the lazy dog"
    rev, revErr := Reverse(input)
    doubleRev, doubleRevErr := Reverse(rev)
    fmt.Printf("original: %q\n", input)
    fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
    fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
}
```

```go
import (
    "errors"
    "fmt"
    "unicode/utf8"
)

func FuzzReverse(f *testing.F) {
    testcases := []string {"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc)  // Use f.Add to provide a seed corpus
    }
    f.Fuzz(func(t *testing.T, orig string) {
        rev, err1 := Reverse(orig)
        if err1 != nil {
            return
        }
        doubleRev, err2 := Reverse(rev)
        if err2 != nil {
            // 除了使用 return，你还可以调用 t.Skip() 来跳过当前的测试输入，继续下一轮测试输入。
             return
        }
        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
```

```go

```
