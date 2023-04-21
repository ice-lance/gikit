/***
 * @Author       : ICE
 * @Date         : 2023-04-19 18:03:42
 * @LastEditTime : 2023-04-20 18:18:58
 * @LastEditors  : ICE
 * @Copyright (c) 2023 ICE, All Rights Reserved.
 * @Description  : 日志处理类
 */
package log

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLog(level string, fields []zapcore.Field) {

	var encoder zapcore.Encoder
	zapEncoderConfig := zap.NewProductionEncoderConfig()
	zapEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	zapEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapEncoderConfig.CallerKey = "line" // 支持中文
	zapEncoderConfig.MessageKey = "msg"
	zapEncoderConfig.LevelKey = "level"
	zapEncoderConfig.TimeKey = "time"

	encoder = zapcore.NewJSONEncoder(zapEncoderConfig)

	// 设置级别
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}

	// 自定义级别展示
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= logLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= logLevel
	})

	// 获取 io.Writer
	infoWriter := getWriter("./logs/info.log")
	warnWriter := getWriter("./logs/error.log")

	// 指向输出
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel), //打印到控制台
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel), zap.Fields(fields...))

}

func getWriter(filename string) io.Writer {
	// filename指向最新的日志
	// 保存 30天,每小时一个文件
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Debug(format string, v ...interface{}) {
	logger.Sugar().Debugf(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Sugar().Infof(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Sugar().Warnf(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Sugar().Errorf(format, v...)
}

func Panic(format string, v ...interface{}) {
	logger.Sugar().Panicf(format, v...)
}
