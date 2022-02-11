---
date: 2020-11-09T22:48:19+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "regexp - 正则表达式"  # 文章标题
url:  "posts/go/libraries/standard/regexp"  # 设置网页链接，默认使用文件名
tags: [ "go", "regexp" ]  # 自定义标签
series: [ "Go 学习笔记" ]  # 文章主题/文章系列
categories: [ "学习笔记" ]  # 分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

- [正则表达式语法规则](#正则表达式语法规则)
		- [1) 字符](#1-字符)
		- [2) 数量词（用在字符或 (...) 之后）](#2-数量词用在字符或--之后)
		- [3) 边界匹配](#3-边界匹配)
		- [4) 逻辑、分组](#4-逻辑分组)
		- [5) 特殊构造（不作为分组）](#5-特殊构造不作为分组)
- [Regexp 标准库示例](#regexp-标准库示例)
	- [匹配指定类型的字符串](#匹配指定类型的字符串)
	- [匹配 a 和 c 中间包含一个数字的字符串](#匹配-a-和-c-中间包含一个数字的字符串)
	- [使用 \d 来匹配 a 和 c 中间包含一个数字的字符串](#使用-d-来匹配-a-和-c-中间包含一个数字的字符串)
	- [匹配字符串中的小数](#匹配字符串中的小数)
	- [匹配、查找、替换功能](#匹配查找替换功能)

## 正则表达式语法规则

#### 1) 字符

| 语法 | 说明 | 表达式示例 | 匹配结果 |
| ---- | ---- | ---- | ---- |
| 一般字符 | 匹配自身 | abc | abc |
| . | 匹配任意除换行符"\n"外的字符， 在 DOTALL 模式中也能匹配换行符 | a.c | abc |
| \ | 转义字符，使后一个字符改变原来的意思； 如果字符串中有字符 * 需要匹配，可以使用 \* 或者字符集［*]。 | a\.c a\\c | a.c a\c |
| [...] | 字符集（字符类），对应的位置可以是字符集中任意字符。 字符集中的字符可以逐个列出，也可以给出范围，如 [abc] 或 [a-c]， 第一个字符如果是 ^ 则表示取反，如 [^abc] 表示除了abc之外的其他字符。 | a[bcd]e | abe 或 ace 或 ade |
| \d | 数字：[0-9] | a\dc | a1c |
| \D | 非数字：[^\d] | a\Dc | abc |
| \s | 空白字符：[<空格>\t\r\n\f\v] | a\sc | a c |
| \S | 非空白字符：[^\s] | a\Sc | abc |
| \w | 单词字符：[A-Za-z0-9] | a\wc | abc |
| \W | 非单词字符：[^\w] | a\Wc | a c |

#### 2) 数量词（用在字符或 (...) 之后）

| 语法 | 说明 | 表达式示例 | 匹配结果 |
| ---- | ---- | ---- | ---- |
| * | 匹配前一个字符 0 或无限次 | abc* | ab 或 abccc |
| + | 匹配前一个字符 1 次或无限次 | abc+ | abc 或 abccc |
| ? | 匹配前一个字符 0 次或 1 次 | abc? | ab 或 abc |
| {m} | 匹配前一个字符 m 次 | ab{2}c | abbc |
| {m,n} | 匹配前一个字符 m 至 n 次，m 和 n 可以省略，若省略 m，则匹配 0 至 n 次； 若省略 n，则匹配 m 至无限次 | ab{1,2}c | abc 或 abbc |

#### 3) 边界匹配

| 语法 | 说明 | 表达式示例 | 匹配结果 |
| ---- | ------- | ---- | ---- |
| ^ | 匹配字符串开头，在多行模式中匹配每一行的开头 | ^abc | abc |
| $ | 匹配字符串末尾，在多行模式中匹配每一行的末尾 | abc$ | abc |
| \A | 仅匹配字符串开头 | \Aabc | abc |
| \Z | 仅匹配字符串末尾 | abc\Z | abc |
| \b | 匹配 \w 和 \W 之间 | a\b!bc | a!bc |
| \B | [^\b] | a\Bbc | abc |

#### 4) 逻辑、分组

| 语法 | 说明 | 表达式示例 | 匹配结果 |
| ---- | ---- | ---- | ---- |
| \| | \| 代表左右表达式任意匹配一个，优先匹配左边的表达式 | abc\|def | abc 或 def |
| (...) | 括起来的表达式将作为分组，分组将作为一个整体，可以后接数量词 | (abc){2} | abcabc |
| (?P<name>...) | 分组，功能与 (...) 相同，但会指定一个额外的别名 | (?P<id>abc){2} | abcabc |
| \<number> | 引用编号为 <number> 的分组匹配到的字符串 | (\d)abc\1 | 1abe1 或 5abc5 |
| (?P=name) | 引用别名为 <name> 的分组匹配到的字符串 | (?P<id>\d)abc(?P=id) | 1abe1 或 5abc5 |

#### 5) 特殊构造（不作为分组）

| 语法 | 说明 | 表达式示例 | 匹配结果 |
| ---- | ---- | ---- | ---- |
| (?:...) | (…) 的不分组版本，用于使用 "\|" 或后接数量词 | (?:abc){2} | abcabc |
| (?iLmsux) | iLmsux 中的每个字符代表一种匹配模式，只能用在正则表达式的开头，可选多个 | (?i)abc | AbC |
| (?#...) | # 后的内容将作为注释被忽略。 | abc(?#comment)123 | abc123 |
| (?=...) | 之后的字符串内容需要匹配表达式才能成功匹配 | a(?=\d) | 后面是数字的 a |
| (?!...) | 之后的字符串内容需要不匹配表达式才能成功匹配 | a(?!\d) | 后面不是数字的 a |
| (?<=...) | 之前的字符串内容需要匹配表达式才能成功匹配 | (?<=\d)a | 前面是数字的a |
| (?<!...) | 之前的字符串内容需要不匹配表达式才能成功匹配 | (?<!\d)a | 前面不是数字的a |

## Regexp 标准库示例

### 匹配指定类型的字符串

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	buf := "abc azc a7c aac 888 a9c  tac"
	reg := regexp.MustCompile(`a.c`)
	if reg == nil {
		return
	}
	result := reg.FindAllStringSubmatch(buf, -1)
	fmt.Println("result = ", result)
}
```

```
result =  [[abc] [azc] [a7c] [aac] [a9c]]
```

### 匹配 a 和 c 中间包含一个数字的字符串

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	buf := "abc azc a7c aac 888 a9c  tac"
	reg := regexp.MustCompile(`a[0-9]c`)
	if reg == nil {
		return
	}
	result := reg.FindAllStringSubmatch(buf, -1)
	fmt.Println("result = ", result)
}
```

```
result = [[a7c] [a9c]]
```

### 使用 \d 来匹配 a 和 c 中间包含一个数字的字符串

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	buf := "abc azc a7c aac 888 a9c  tac"
	reg := regexp.MustCompile(`a\dc`)
	if reg == nil {
		return
	}
	result := reg.FindAllStringSubmatch(buf, -1)
	fmt.Println("result = ", result)
}
```

```
result = [[a7c] [a9c]]
```

### 匹配字符串中的小数

```go
package main

import (
	"fmt"
	"regexp"
)

func main() {
	buf := "43.14 567 agsdg 1.23 7. 8.9 1sdljgl 6.66 7.8   "
	reg := regexp.MustCompile(`\d+\.\d+`)
	if reg == nil {
		return
	}
	result := reg.FindAllString(buf, -1)
	//result := reg.FindAllStringSubmatch(buf, -1)
	fmt.Println("result = ", result)
}
```

```
result = [[43.14] [1.23] [8.9] [6.66] [7.8]]
```

### 匹配、查找、替换功能

```go
package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	pat := "[0-9]+.[0-9]+"
	f := func(s string) string {
		v, _ := strconv.ParseFloat(s, 32)
		return strconv.FormatFloat(v*2, 'f', 2, 32)
	}
	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
		fmt.Println("Match Found!")
	}
	re, _ := regexp.Compile(pat)
	str := re.ReplaceAllString(searchIn, "##.#")
	fmt.Println(str)
	str2 := re.ReplaceAllStringFunc(searchIn, f)
	fmt.Println(str2)
}
```

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	searchIn := "control_circuit_breaker_signal必须是[0 1]中的一个"

	searchIn = regexp.MustCompile(`([\w(),<>\-{}/\+\.\[\]]+)`).ReplaceAllString(searchIn, " $1 ")
	searchIn = regexp.MustCompile(` {2,}`).ReplaceAllString(searchIn, " ")
	searchIn = strings.TrimSpace(searchIn)

	fmt.Println(searchIn)
}
```
