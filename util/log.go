package util

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	hook := lumberjack.Logger{
		Filename:   "log/wechat.log",
		MaxSize:    1024,
		MaxBackups: 100,
		MaxAge:     10,
		Compress:   true,
	}
	level := new(zapcore.Level)
	if err := level.Set("debug"); err != nil {
		*level = zap.DebugLevel
	}
	w := zapcore.AddSync(&hook)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		*level,
	)
	Logger = zap.New(core, zap.AddCaller())
}
