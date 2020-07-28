# 字符串与其他类型相互转换

strconv 标准库实现了基本数据类型和其字符串表示的相互转换。

## string 与 int 类型转换

**Atoi()** - 将字符串类型的整数转换为 int 类型

```go
// array to int
func Atoi(s string) (i int, err error)
```

**Itoa()** - 将 int 类型数据转换为对应的字符串表示

```go
// int to array
func Itoa(i int) string
```

> C 语言遗留问题，C 语言中字符串是数组表示的。

## Parse 系列函数

转换字符串为给定类型的值。

- ParseBool()

```go
func ParseBool(str string) (value bool, err error)
```

- ParseInt()

```go
func ParseInt(s string, base int, bitSize int) (i int64, err error)
```

- ParseUint() - 无符号

```go
func ParseUint(s string, base int, bitSize int) (n uint64, err error)
```

- ParseFloat()

```go
func ParseFloat(s string, bitSize int) (f float64, err error)
```

## Format 系列函数

将给定类型数据格式化为字符串类型数据。

- FormatBool()

```go
func FormatBool(b bool) string
```

- FormatInt()

```go
func FormatInt(i int64, base int) string
```

- FormatUint() - 无符号

```go
func FormatUint(i uint64, base int) string
```

base 表示进制。

- FormatFloat()

```go
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```

bitSize 表示 f 的来源类型（float32、float64），会据此进行舍入。

## 其他

- isPrint() - 返回一个字符是否是可打印的。

```go
func IsPrint(r rune) bool
```

- CanBackquote() - 返回字符串是否可以不被修改的表示为一个单行的、没有空格符号之外控制字符的反引号字符串。

```go
func CanBackquote(s string) bool
```
