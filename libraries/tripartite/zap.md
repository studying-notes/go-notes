# Zap 日志库

- [Zap 日志库](#zap-日志库)
	- [介绍](#介绍)
	- [Zap Logger](#zap-logger)
		- [Logger](#logger)
		- [Sugared Logger](#sugared-logger)
	- [定制 Logger](#定制-logger)
		- [将日志写入文件而不是终端](#将日志写入文件而不是终端)
		- [将 JSON Encoder 更改为普通的 Log Encoder](#将-json-encoder-更改为普通的-log-encoder)
		- [更改时间编码/添加调用者详细信息](#更改时间编码添加调用者详细信息)
	- [用 Lumberjack 根据文件大小进行日志切割归档](#用-lumberjack-根据文件大小进行日志切割归档)
		- [加入 Lumberjack](#加入-lumberjack)
		- [完整示例](#完整示例)
	- [用 file-foratelogs 根据文件大小进行日志切割归档](#用-file-foratelogs-根据文件大小进行日志切割归档)

## 介绍

一个优秀的日志记录器能够提供下面这些功能：

- 能够将事件记录到文件中，而不只是标准输出；
- 日志切割，即能够根据文件大小、时间或间隔等来切割日志文件；
- 支持不同的日志级别。例如 `INFO`，`DEBUG`，`ERROR` 等；
- 能够打印基本信息，如调用文件、函数名、行号、时间等。

[Zap](https://github.com/uber-go/zap) 是非常快的、结构化的，分日志级别的 Go 语言日志库。同时提供了结构化日志记录和 `Printf` 风格的日志记录，性能比类似的结构化日志库更好。

```shell
go get -u go.uber.org/zap
```

## Zap Logger

Zap 提供了两种类型的日志记录器：`Sugared Logger` 和 `Logger`。

在性能很好但不是很关键的上下文中，可以用 `SugaredLogger`。它比其他结构化日志记录快 4-10 倍，并且支持结构化和 `Printf` 风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，建议用` Logger`。它甚至比 `SugaredLogger` 更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。

### Logger

- 通过调用 `zap.NewProduction()`/`zap.NewDevelopment()` 或者 `zap.Example()` 创建一个 Logger。
- 上面的每一个函数都将创建一个Logger。唯一的区别在于它将记录的信息不同。例如 `Production Logger` 默认记录调用函数信息、日期和时间等。
- 通过 Logger 调用 Info/Error 等。
- 默认情况下日志都会打印到标准输出。

```go
package main

import (
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.Logger

func InitLogger() {
	logger, _ = zap.NewProduction()
}

func HttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching url",
			zap.String("url", url), zap.Error(err))
	} else {
		logger.Info("Success",
			zap.String("status", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer logger.Sync()
	HttpGet("baidu.com")
	HttpGet("http://www.baidu.com")
}
```

日志记录器方法的语法是这样的：

```go
func (log *Logger)Method(msg string, fields ...Field) 
```

其中 `Method` 是一个可变参数函数，可以是 Info/Error/Debug/Panic 等。每个方法都接受一个消息字符串和任意数量的 `zapcore.Field` 参数。每个 `zapcore.Field` 其实就是一组键值对参数。

运行程序，输出如下结果：

```json
{"level":"error","ts":1595375411.0869894,"caller":"log/zap.go:17","msg":"Error fetching url","url":"baidu.com","error":"Get baidu.com: unsupported protocol scheme \"\"","stacktrace":"main.HttpGet\n\tD:/go-notes/io/log/zap.go:17\nmain.main\n\tD:/go-notes/io/log/zap.go:30\nruntime.main\n\tC:/Developer/Go/src/runtime/proc.go:203"}

{"level":"info","ts":1595375412.1307695,"caller":"log/zap.go:20","msg":"Success","status":"200
 OK","url":"http://www.baidu.com"}
```

### Sugared Logger

用 Sugared Logger 来实现相同的功能。

- 大部分的实现基本都相同；
- 惟一的区别是，我们通过调用主 Logger 的 `. Sugar()` 方法来获取一个 `SugaredLogger`；
- 然后使用 `SugaredLogger` 以 `printf` 格式记录语句。

下面是修改过后使用 `SugaredLogger` 代替 `Logger` 的代码：

```go
package main

import (
	"go.uber.org/zap"
	"net/http"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func InitLogger() {
	logger, _ = zap.NewProduction()
	sugaredLogger = logger.Sugar()
}

func HttpGet(url string) {
	sugaredLogger.Debugf("Try to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugaredLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugaredLogger.Infof("Success! status = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer logger.Sync()
	HttpGet("baidu.com")
	HttpGet("http://www.baidu.com")
}
```

```json
{"level":"error","ts":1595375873.972989,"caller":"log/zap.go:20","msg":"Error fetching URL baidu.com : Error = Get baidu.com: unsupported protocol scheme \"\"","stacktrace":"main.HttpGet\n\tD:/go-notes/io/log/zap.go:20\nmain.main\n\tD:/go-notes/io/log/zap.go:30\nruntime.main\n\tC:/Developer/Go/src/runtime/proc.go:203"}

{"level":"info","ts":1595375874.0140352,"caller":"log/zap.go:22","msg":"Success! status = 200 OK
for URL http://www.baidu.com"}
```

到目前为止这两个 Logger 都打印输出 `JSON` 结构格式。

## 定制 Logger

### 将日志写入文件而不是终端

第一个更改是把日志写入文件，而不是打印到应用程序控制台。

- 用 `zap.New()` 方法来手动传递所有配置，而不是使用像 `zap.NewProduction()` 这样的预置方法来创建 Logger。

```go
func New(core zapcore.Core, options ...Option) *Logger
```

`zapcore.Core` 需要三个配置：`Encoder`，`WriteSyncer`，`LogLevel`。

1. **Encoder**：编码器，指定如何写入日志。我们将使用开箱即用的 `NewJSONEncoder()`，并使用预先设置的 `ProductionEncoderConfig()`。

```go
zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
```

2. **WriterSyncer** ：指定日志将写到哪里去。我们使用 `zapcore.AddSync()` 函数并且将打开的文件句柄传进去。

```go
file, _ := os.Create("./dev.log")
writeSyncer := zapcore.AddSync(file)
```

3. **Log Level**：哪种级别的日志将被写入。

我们将修改上述部分中的 Logger 代码，并重写 `InitLogger()` 方法。

```go
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var sugaredLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugaredLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("dev.log")
	return zapcore.AddSync(file)
}

func HttpGet(url string) {
	sugaredLogger.Debugf("Try to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugaredLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugaredLogger.Infof("Success! status = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
    InitLogger()
    defer sugaredLogger.Sync()
	HttpGet("baidu.com")
	HttpGet("http://www.baidu.com")
}
```

当使用这些修改过的配置调用上述部分的 `main()` 函数时，以下输出将打印在文件 `dev.log` 中。

```json
{"level":"debug","ts":1595376445.2954438,"msg":"Try to hit GET request for baidu.com"}
{"level":"error","ts":1595376445.2964363,"msg":"Error fetching URL baidu.com : Error = Get baidu.com: unsupported protocol scheme \"\""}
{"level":"debug","ts":1595376445.2964363,"msg":"Try to hit GET request for http://www.baidu.com"}
{"level":"info","ts":1595376445.3403177,"msg":"Success! status = 200 OK for URL http://www.baidu.com"}

```

### 将 JSON Encoder 更改为普通的 Log Encoder

现在，我们希望将编码器从 `JSON Encoder` 更改为普通 Encoder。为此，我们需要将 `NewJSONEncoder()` 更改为 `NewConsoleEncoder()`。

```go
return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
```

输出：

```shell
1.5953766628413453e+09	debug	Try to hit GET request for baidu.com
1.5953766628423553e+09	error	Error fetching URL baidu.com : Error = Get baidu.com: unsupported protocol scheme ""
1.5953766628423553e+09	debug	Try to hit GET request for http://www.baidu.com
1.595376662882778e+09	info	Success! status = 200 OK for URL http://www.baidu.com
```

### 更改时间编码/添加调用者详细信息

第一件事是覆盖默认的 `ProductionConfig()`，进行以下更改:

- 修改时间编码器
- 在日志文件中使用大写字母记录日志级别

```go
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
```

然后修改 Zap Logger 代码，添加将调用函数信息记录到日志中的功能。在 `zap.New()` 函数中添加一个 `Option`。

```go
logger := zap.New(core, zap.AddCaller())
```

最终输出：

```shell
2020-07-22T08:17:14.981+0800	DEBUG	log/zap.go:38	Try to hit GET request for baidu.com
2020-07-22T08:17:14.997+0800	ERROR	log/zap.go:41	Error fetching URL baidu.com : Error = Get baidu.com: unsupported protocol scheme ""
2020-07-22T08:17:14.997+0800	DEBUG	log/zap.go:38	Try to hit GET request for http://www.baidu.com
2020-07-22T08:17:15.049+0800	INFO	log/zap.go:43	Success! status = 200 OK for URL http://www.baidu.com
```

## 用 Lumberjack 根据文件大小进行日志切割归档

> Zap 本身不支持切割归档日志文件

```shell
go get -u github.com/natefinch/lumberjack
```

### 加入 Lumberjack

修改 `WriteSyncer` 代码，按照下面的代码修改 `getLogWriter()` 函数：

```go
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "dev.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

Lumberjack Logger 属性:

- `Filename`: 日志文件的位置
- `MaxSize`：在进行切割之前，日志文件的最大大小（MB）
- `MaxBackups`：保留旧文件的最大个数
- `MaxAges`：保留旧文件的最大天数
- `Compress`：是否压缩/归档旧文件

### 完整示例

```go
package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var sugaredLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	//logger := zap.New(core)
	sugaredLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)

	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "dev.log",
		MaxSize:    1,
		MaxAge:     5,
		MaxBackups: 30,
		LocalTime:  false,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func HttpGet(url string) {
	sugaredLogger.Debugf("Try to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugaredLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugaredLogger.Infof("Success! status = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer sugaredLogger.Sync()
	HttpGet("baidu.com")
	HttpGet("http://www.baidu.com")
}
```

## 用 file-foratelogs 根据文件大小进行日志切割归档

```shell
go get github.com/lestrrat-go/file-rotatelogs
```

```go
package main

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

var Logger *zap.Logger

func InitLogger() {
	// 判断日志等级
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel
	})

	// 为不同级别的 Logger 创建不同的日志文件
	// 后缀和时间自动添加，只需文件名
	infoWriter := getWriteSyncer("info")
	warnWriter := getWriteSyncer("error")

	encoder := getEncoder()

	// 创建具体的 Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
	)

	// zap.AddCaller() 显示打日志点的文件名和行数
	Logger = zap.New(core, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	// 设置基本日志格式
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// 这个时间字符串似乎是不能乱改的
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	return encoder
}

func getWriteSyncer(filename string) zapcore.WriteSyncer {
	rotateLogger, _ := rotatelogs.New(
		filename+"%Y%m%d%H"+".log",
		// 指向最新日志的链接，Windows 不支持
		//rotatelogs.WithLinkName(filename+".log"),
		// 保存 7 天内的日志
		rotatelogs.WithMaxAge(time.Hour*24*7),
		// 每 1 小时（整点）分割一次日志
		rotatelogs.WithRotationTime(time.Hour),
	)
	return zapcore.AddSync(rotateLogger)
}

func HttpGet(url string) {
	Logger.Debug("Try to hit GET request for" + url)
	resp, err := http.Get(url)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error fetching URL %s : Error = %s", url, err))
	} else {
		Logger.Info(fmt.Sprintf("Success! status = %s for URL %s", resp.Status, url))
		resp.Body.Close()
	}
}

func main() {
	InitLogger()
	defer Logger.Sync()
	HttpGet("baidu.com")
	HttpGet("http://www.baidu.com")
}
```
