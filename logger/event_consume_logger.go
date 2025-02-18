package logger

import (
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

var _event_consume_logger *zap.Logger

func InitEventConsumeLogger(fileName string) {
	eventConsumeLogger := &lumberjack.Logger{
		Filename:   filepath.Join(mviper.GetString("VirtualRegistryManagement.logger.event_consume_log.path"), fileName),
		MaxSize:    mviper.GetInt("VirtualRegistryManagement.logger.event_consume_log.max_size"),    // MB
		MaxBackups: mviper.GetInt("VirtualRegistryManagement.logger.event_consume_log.max_backups"), // how many backup for log file.
		MaxAge:     mviper.GetInt("VirtualRegistryManagement.logger.event_consume_log.max_age"),
		Compress:   mviper.GetBool("VirtualRegistryManagement.logger.event_consume_log.compress"),
	}
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:     "T",
		LevelKey:    "L",
		MessageKey:  "M",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})
	fileWriter := zapcore.AddSync(eventConsumeLogger)
	loggerCore := zapcore.NewCore(encoder, fileWriter, zapcore.InfoLevel)
	core := zapcore.NewTee(loggerCore)
	_event_consume_logger = zap.New(core, zap.AddCaller(), zap.Development())
}

func Use() *zap.Logger {
	return _event_consume_logger
}
