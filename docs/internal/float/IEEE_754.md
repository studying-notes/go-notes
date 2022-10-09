---
date: 2022-10-02T09:04:58+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "IEEE-754 浮点数标准"  # 文章标题
url:  "posts/go/docs/internal/float/IEEE_754"  # 设置网页永久链接
tags: [ "Go", "ieee-754" ]  # 标签
categories: [ "Go 学习笔记" ]  # 分类

toc: true  # 目录
draft: true  # 草稿
---

IEEE-754 规范使用以 2 为底数的指数表示小数，这和使用以 10 为底数的指数表示法（即科学计数法）非常类似。

表 2-1 给出了几个例子，如 0.085 可以用指数的形式表示为 1.36×2^-4，其中 1.36 为系数，2 为底数， -4 为指数。

![](../../../assets/images/docs/internal/float/IEEE_754/表2-1%20数字的表示方法示例.png)

IEEE-754 的浮点数存在多种精度。很显然，更多的存储位数可以表达更大的数或更高的精度。在高级语言中一般存在两种精度的浮点数，即大部分硬件浮点数单元支持的 32 位的单精度浮点数与 64 位的双精度浮点数。

如表 2-2 所示，两种精度的浮点数具有不同的格式。

![](../../../assets/images/docs/internal/float/IEEE_754/表2-2%20单精度与双精度浮点数格式.png)

其中，最开头的 1 位为符号位，1 代表负数，0 代表正数。符号位之后为指数位，单精度为 8 位，双精度为 11 位。指数位存储了指数加上偏移量的值，偏移量是为了表达负数而设计的。例如当指数为 -4 时，实际存储的值为 -4+127 = 123。剩下的是小数位，小数位存储系数中小数位的准确值或最接近的值，是 0 到 1 之间的数。小数位占用的位数最多，直接决定了精度的大小。

以数字 0.085 为例，单精度下的浮点数表示如表 2-3 所示。

![](../../../assets/images/docs/internal/float/IEEE_754/表2-3%20数字0.085的单精度浮点数表示.png)

## 小数部分计算

小数部分的计算是最复杂的，其存储的可能是系数的近似值而不是准确值。小数位的每一位代表的都是 2 的幂，并且指数依次减少 1。以 0.085 的浮点表示法中系数的小数部分（0.36）为例，对应的二进制数为 010 1110 0001 0100 0111 1011，其计算步骤如表 2-4 所示，存储的数值接近 0.36。

| 位 | 对应的整数 | 转化为分数 | 转化为十进制小数 | 各位的总和 |
| ---- | -------- | -------- | -------- | -------- |
| 2 | 4 | 1/4 | 0.25 | 0.25 |
| 4 | 16 | 1/16 | 0.0625 | 0.3125 |
| 5 | 32 | 1/32 | 0.03125 | 0.34375 |
| 6 | 64 | 1/64 | 0.015625 | 0.359375 |
| 11 | 2048 | 1/2048 | 0.00048828125 | 0.35986328125 |
| 13 | 8192 | 1/8192 | 0.0001220703125 | 0.3599853515625 |
| 17 | 131072 | 1/131072 | 0.00000762939453 | 0.35999298095703 |
| 18 | 262144 | 1/262144 | 0.00000381469727 | 0.3599967956543 |
| 19 | 524288 | 1/524288 | 0.00000190734863 | 0.35999870300293 |
| 20 | 1048576 | 1/1048576 | 0.00000095367432 | 0.35999965667725 |
| 22 | 4194304 | 1/4194304 | 0.00000023841858 | 0.35999989509583 |
| 23 | 8388608 | 1/8388608 | 0.00000011920929 | 0.36000001430512 |

那么小数位又是如何计算出来的呢？以数字 0.085 为例，可以使用“乘 2 取整法”将该十进制小数转化为二进制小数，即

```
0.085（十进制）
=0.00010101110000101000111101011100001010001111010111000011（二进制）
=1.0101110000101000111101011100001010001111010111000011×2^-4
```

由于小数位只有 23 位，因此四舍五入后为 010 1110 0001 0100 0111 1011，这就是最终浮点数的小数部分。

### 显示浮点数格式

Go 语言标准库的 math 包提供了许多有用的计算函数，其中，Float32 可以以字符串的形式打印出单精度浮点数的二进制值。

下例中的 Go 代码可以输出 0.085 的浮点数表示中的符号位、指数位与小数位。

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	var number float32 = 0.085
	fmt.Printf("Number: %f\n", number)

	bits := math.Float32bits(number)
	fmt.Printf("Bits: %b\n", bits)
	fmt.Printf("Binary: %.32b\n", bits)

	binary := fmt.Sprintf("%.32b", bits)
	fmt.Printf(
		"Pattern: %s | %s %s | %s %s %s %s %s %s\n",
		binary[0:1],
		binary[1:5],
		binary[5:9],
		binary[9:12],
		binary[12:16],
		binary[16:20],
		binary[20:24],
		binary[24:28],
		binary[28:32],
	)
}
```

```
Number: 0.085000
Bits: 111101101011100001010001111011
Binary: 00111101101011100001010001111011
Pattern: 0 | 0111 1011 | 010 1110 0001 0100 0111 1011
```

为了验证之前理论的正确性，可以根据二进制值反向推导出其所表示的原始十进制值 0.085。思路是将符号位、指数位、小数位分别提取出来，将小数部分中每个为 1 的 bit 位都转化为对应的十进制小数，并求和。

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	var number float32 = 0.085
	bits := math.Float32bits(number)

	bias := uint32(127)
	signBit := bits >> 31           // 符号位
	exponent := (bits >> 23) & 0xFF // 指数部分
	exponent = exponent - bias      // 指数部分偏移
	mantissa := bits & 0x7FFFFF     // 小数部分

	fmt.Printf("signBit: %b\n", signBit)
	fmt.Printf("exponent: %08b, %d\n", uint8(exponent), int8(exponent))
	fmt.Printf("mantissa: %023b\n", mantissa)

	var value float32

	// 还原小数部分
	for i := 0; i < 23; i++ {
		if mantissa&(1<<i) != 0 {
			value += float32(1 / math.Pow(2, float64(23-i)))
		}
	}

	fmt.Printf("mantissa: %f\n", value)

	value = (1 + value) * float32(math.Pow(2, float64(int8(exponent))))

	if signBit == 1 {
		value = -value
	}

	fmt.Printf("value: %f\n", value)
}
```

符号位、指数位、小数位，以及最终结果输出如下，验证了之前的理论。

```
signBit: 0
exponent: 11111100, -4
mantissa: 01011100001010001111011
mantissa: 0.360000
value: 0.085000
```

## 判断浮点数为整数

判断浮点数为整数的重要思路是指数能够弥补小数部分（即指数的值大于或等于小数的位数）。

例如，在十进制数中，1.23×10^2 是整数，而 1.234×10^2 不是整数，因为指数 2 不能弥补 3 个小数位。以 2 为底数的浮点数的判断思路类似。

```go
func isInt(bits uint32) bool {
	exponent := int8((bits>>23)&0xFF) - 127 // 127 is the bias for float32
	mantissa := bits & 0x7FFFFF             // 23 bits
	return exponent >= 0 &&                 // exponent is positive
		mantissa&(0x7FFFFF>>uint(exponent)) == 0 // check if mantissa is an integer
}

func Example_isInt() {
	fmt.Println(isInt(math.Float32bits(0.1)))
	fmt.Println(isInt(math.Float32bits(1)))
	fmt.Println(isInt(math.Float32bits(1.1)))
	fmt.Println(isInt(math.Float32bits(10.1)))
	fmt.Println(isInt(math.Float32bits(10.0)))
	fmt.Println(isInt(math.Float32bits(10.0001)))
	fmt.Println(isInt(math.Float32bits(0.10101)))
	fmt.Println(isInt(math.Float32bits(12345)))

	// Output:
	// false
	// true
	// false
	// false
	// true
	// false
	// false
	// true
}
```

要保证浮点数格式中实际存储的数为整数，一个必要条件就是浮点数格式中指数位的值大于 127。指数位的值为 127 代表指数为 0，如果指数位的值大于 127，则代表指数大于 0，反之则代表指数小于 0。

只有当指数可以弥补小数部分的时候，才是一个整数。例如，数字 234523 的指数的值是 144-127 = 17，代表其不能弥补最后 23-17 = 6 位的小数，即当最后 6 位不全为 0 时，数字 234523 一定为小数。但由于数字 234523 最后 6 位刚好都为 0，所以可以判断它是整数。

## 常规数与非常规数

在 IEEE-754 中指数位有一个偏移量，偏移量是为了表达负数而设计的。

比如单精度中的 0.085，实际的指数是 -4，在点数格式中指数位存储的是数字 123。所以可以看出，浮点数指数位表达的负数值始终是有下限的。单精度浮点数指数值的下限就是 -126。如果比这个数还要小，例如 -127，那么应该表达为 0.1×2-126，这时的系数小于 1。

我们把系数小于 1 的数叫作非常规数（Denormal Number），把系数在 1 到 2 之间的数叫作常规数（Normal Number）。

## NaN与Inf

在 Go 语言中有正无穷(+Inf)与负无穷(-Inf)两类异常的值，例如正无穷 1/0。

NaN 代表异常或无效的数字，例如 0/0 或者 Sqrt(-1)。

在 IEEE754 浮点数标准中，在正常情况下，不可能所有的指数位都为 1 或者都为 0。

例如，Float32 的最大值其实是 0|1111 1110|111 1111 1111 1111 1111 1111。

当所有的指数位都为 0 时代表 0，当所有的指数位都为 1 时代表 -1。

在 IEEE-754 标准中，NaN 分为 sNAN 与 qNAN。qNAN 代表出现了无效或异常的结果，sNAN 代表发生了无效的操作，例如将字符串转化为浮点数。qNAN 的指数位全为 1，且小数位的第一位为 1 ； sNAN 的指数位全为 1，但是小数位的第一位为 0。

用 math.NaN 函数可以生成一个 NaN，对 NAN 的任何操作都会返回 NAN。另外，对 NAN 的任何大小比较都会返回 false。例如：

```go
func main() {
	nan := math.NaN()
	fmt.Println(nan == nan, .nan < nan, nan > nan, nan <= nan, nan >= nan)
	// false false false false false
}
```

有些时候需要判断浮点数是否为 NaN 或者 Inf，这需要借助 Math.IsNaN 和 Math.IsInf 函数。

其判断条件很简单，在 IEEE-754 标准中，`NaN!=NaN` 会返回 true，Go 语言编译器在判断浮点数时，浮点数的比较会被编译成 `UCOMISD` 或 `COMISD` CPU 指令，该指令会判断和处理 NaN 等异常情况从而实现当 `NaN!=NaN` 时返回 true。可以通过判断浮点数是否在有效的范围内来检查其是否为 Inf。

浮点数的最大和最小值的常量在 `src/math/const.go` 中定义。

```go
// Floating-point limit values.
// Max is the largest finite value representable by the type.
// SmallestNonzero is the smallest positive, non-zero value representable by the type.
const (
	MaxFloat32             = 0x1p127 * (1 + (1 - 0x1p-23)) // 3.40282346638528859811704183484516925440e+38
	SmallestNonzeroFloat32 = 0x1p-126 * 0x1p-23            // 1.401298464324817070923729583289916131280e-45

	MaxFloat64             = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
	SmallestNonzeroFloat64 = 0x1p-1022 * 0x1p-52            // 4.9406564584124654417656879286822137236505980e-324
)
```

`src/math/bits.go`

```go
// IsNaN reports whether f is an IEEE 754 “not-a-number” value.
func IsNaN(f float64) (is bool) {
	// IEEE 754 says that only NaNs satisfy f != f.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return uint32(x>>shift)&mask == mask && x != uvinf && x != uvneginf
	return f != f
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf(f float64, sign int) bool {
	// Test for infinity by comparing against maximum float.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return sign >= 0 && x == uvinf || sign <= 0 && x == uvneginf;
	return sign >= 0 && f > MaxFloat64 || sign <= 0 && f < -MaxFloat64
}
```

```go

```
