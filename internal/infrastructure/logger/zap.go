package logger

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(config *viper.Viper) *zap.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.GetString("LOG_FILE_PATH"),
		MaxSize:    config.GetInt("LOG_MAX_SIZE"),
		MaxBackups: config.GetInt("LOG_MAX_BACKUPS"),
		MaxAge:     config.GetInt("LOG_MAX_AGE"),
		Compress:   config.GetBool("LOG_COMPRESS"),
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(lumberjackLogger), zap.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)

	return zap.New(core)
}
