# Logrus 日志库

Logrus 是 Go 的结构化 Logger，与 Log 标准库完全兼容。

它有以下特点：

- 完全兼容标准日志库，拥有七种日志级别：`Trace`, `Debug`, `Info`, `Warning`, `Error`, `Fatal`, `Panic`；
- 可扩展的 `Hook` 机制，允许使用者通过 `Hook` 的方式将日志分发到任意地方；
- 可选的日志输出格式，内置了两种日志格式 `JSONFormater` 和 `TextFormatter`，还可以自定义日志格式；
- `Field` 机制，通过 `Filed` 机制进行结构化的日志记录；
- 线程安全。

```
go get github.com/sirupsen/logrus
```

## 基本示例

最简单的是 `Package-Level` 导出日志程序：

```go
package main

import (
  log "github.com/sirupsen/logrus"
)

func main() {
  log.WithFields(log.Fields{
    "animal": "walrus",
  }).Info("A walrus appears")
}
```

```go
package main

import (
  "os"
  log "github.com/sirupsen/logrus"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr
  // Can be any io.Writer, see below for File example
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.WarnLevel)
}

func main() {
  log.WithFields(log.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

  log.WithFields(log.Fields{
    "omg":    true,
    "number": 122,
  }).Warn("The group's number increased tremendously!")

  log.WithFields(log.Fields{
    "omg":    true,
    "number": 100,
  }).Fatal("The ice breaks!")

  // A common pattern is to re-use fields between logging statements by re-using
  // the logrus.Entry returned from WithFields()
  contextLogger := log.WithFields(log.Fields{
    "common": "this is a common field",
    "other": "I also should be logged always",
  })

  contextLogger.Info("I'll be logged with common and other field")
  contextLogger.Info("Me too")
}
```

## 进阶示例

记录到多个位置：

```go
package main

import (
  "os"
  "github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func main() {
  // The API for setting attributes is a little different than the package level
  // exported logger. See Godoc.
  log.Out = os.Stdout

  // You could set this to any `io.Writer` such as a file
  // file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  // if err == nil {
  //  log.Out = file
  // } else {
  //  log.Info("Failed to log to file, using default stderr")
  // }

  log.WithFields(logrus.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")
}
```

## 日志级别

Logrus 有七个日志级别：`Trace`, `Debug`, `Info`, `Warning`, `Error`, `Fatal`, `Panic`。

```go
log.Trace("Something very low level.")
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
```

### 设置日志级别

可以在 Logger 上设置日志记录级别，然后它只会记录具有该级别或以上级别任何内容的条目：

```go
log.SetLevel(log.InfoLevel)
```

## 字段

Logrus 鼓励通过日志字段进行谨慎的结构化日志记录，而不是冗长的、不可解析的错误消息。

例如，区别于使用 `log.Fatalf("Failed to send event %s to topic %s with key %d")`，推荐使用如下方式记录更容易发现的内容：

```go
log.WithFields(log.Fields{
  "event": event,
  "topic": topic,
  "key": key,
}).Fatal("Failed to send event")
```

## 默认字段

通常，将一些字段始终附加到应用程序的全部或部分的日志语句中会很有帮助。例如，可能希望始终在请求的上下文中记录 `request_id` 和 `user_ip`。

区别于在每一行日志中写上 `log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})`，可以向下面的示例代码一样创建一个 `logrus.Entry` 去传递这些字段。

```go
requestLogger := log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})
requestLogger.Info("something happened on that request")
requestLogger.Warn("something not great happened")
```

## 日志条目

除了使用 `WithField` 或 `WithFields` 添加的字段外，一些字段会自动添加到所有日志记录事中:

- time：记录日志时的时间戳
- msg：记录的日志信息
- level：记录的日志级别

## Hooks

可以添加日志级别的钩子（Hook）。例如，向异常跟踪服务发送 `Error`、`Fatal` 和 `Panic` 信息到 StatsD 或同时将日志发送到多个位置，例如 syslog。

Logrus 配有内置钩子，在 `init` 中添加这些内置钩子或自定义的钩子：

```go
import (
  log "github.com/sirupsen/logrus"
  "gopkg.in/gemnasium/logrus-airbrake-hook.v2" // the package is named "airbrake"
  logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
  "log/syslog"
)

func init() {

  // Use the Airbrake hook to report errors that have Error severity or above to
  // an exception tracker. You can create custom hooks, see the Hooks section.
  log.AddHook(airbrake.NewHook(123, "xyz", "production"))

  hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
  if err != nil {
    log.Error("Unable to connect to local syslog daemon")
  } else {
    log.AddHook(hook)
  }
}
```

## 格式化

logrus内置以下两种日志格式化程序：

```
logrus.TextFormatter
logrus.JSONFormatter
```

## 记录函数名

如果希望将调用的函数名添加为字段，请通过以下方式设置：


```go
log.SetReportCaller(true)
```

## 线程安全

默认的 logger 在并发写的时候是被 mutex 保护的，比如当同时调用 hook 和写 log 时 mutex 就会被请求，有另外一种情况，文件是以 appending mode 打开的， 此时的并发操作就是安全的，可以用 `logger.SetNoLock()` 来关闭它。

## Gin 框架使用 Logrus

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.JSONFormatter{}
	f, _ := os.Create("gin.log")
	log.Out = f
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Out
	log.Level = logrus.InfoLevel
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		log.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Warn("A group of walrus emerges from the ocean.")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	_ = r.Run()
}
```
