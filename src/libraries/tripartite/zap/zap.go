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
			// 这个时间字符串是不能乱改的
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
