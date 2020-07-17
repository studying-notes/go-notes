# Windows 10 64x 安装 go-sqlcipher

## 预准备

1. 安装 Perl 64bit

有不同的发行版本，随便选一个即可，比如：

```
http://strawberryperl.com/download/5.30.2.1/strawberry-perl-5.30.2.1-64bit.msi
```

安装到默认或者指定路径，然后添加 `bin` 目录到 `PATH` 中。

2. 安装 GCC 64bit

也有很多选择，比如 TDM-GCC-64：

```
https://jmeubank.github.io/tdm-gcc/
```

3. 安装 MSYS

```
http://downloads.sourceforge.net/mingw/MSYS-1.0.11.exe
```

不可缺少，缺了，一些命令无法在 Windows 下无法运行。

4. 安装 Make

无论如何搞定，最后可以提供 `make` 命令即可，`Powershell` 下可以以管理员身份执行以下命令安装：

```powershell
> choco install make
```

另一种：

```powershell
> copy mingw32-make.exe make.exe
```

启动一个终端，运行以下命令，都可以正确显示就表明准备工作就绪了：

```powershell
> perl -v
> gcc -v
> make -v
```

## 编译 OpenSSL 64bit

1. 下载 OpenSSL

```
https://www.openssl.org/source/
```

新版本编译起来更复杂，这里用的旧版本：

```
https://www.openssl.org/source/openssl-1.0.2k.tar.gz
```

2. 解压

```powershell
> tar -xvzf openssl-1.0.2k.tar.gz
```

用 360 解压缩这一类工具也没问题。

3. 编译

```powershell
> cd openssl-1.0.2k
> perl configure mingw64 no-shared no-asm
> make
```

编译过程大概十几分钟。

编译完成后，将 `openssl-1.0.2k` 目录下的以下两个文件 `libcrypto.a`、`libcrypto.pc` 复制到 `TDM-GCC-64\lib` 目录下，然后将 `openssl-1.0.2k\include\openssl` 这个文件夹复制到 `TDM-GCC-64\x86_64-w64-mingw32\include` 下。

## 安装 go-sqlcipher

1. 下载源码

```powershell
> go get github.com/xeodou/go-sqlcipher
```

不出意外，安装报错，但源码已经下载下来了，找到所在目录，一般在：

```
GOPATH\pkg\mod\github.com\xeodou
```

这里我遇到了一个问题，全局 GOPATH 和 Goland 的 GOPATH 设置的不一样，导致 Goland 没有问题而命令行运行一直出错。

2. 修改 `sqlite3_windows.go` 文件

打开一级目录下的  `sqlite3_windows.go` 文件，将开始部分改成以下内容：

```go
// Copyright (C) 2014 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
// +build windows

package sqlite3

/*
#cgo CFLAGS: -I. -fno-stack-check -fno-stack-protector -mno-stack-arg-probe
#cgo windows,386 CFLAGS: -D_USE_32BIT_TIME_T
#cgo LDFLAGS: -lmingwex -lmingw32 -lgdi32
*/
import "C"
```

3. 再次安装

```powershell
> go install -v .
```

完成！
