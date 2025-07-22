package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/Zillaforge/toolkits/mviper"
)

// Init is responsible to initialize logger.
// When the level in configuration(YAML) is set high, logger will be recorded all information. Including Info, Error etc.
// When the level is set low, logger will just record Error information in log file.
func Init(fileName string) {
	mode := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		switch mviper.GetString("VirtualRegistryManagement.logger.system_log.mode") {
		case "debug":
			return lvl >= zapcore.DebugLevel
		case "info":
			return lvl >= zapcore.InfoLevel
		default:
			return lvl >= zapcore.ErrorLevel
		}
	})

	hook := &lumberjack.Logger{
		Filename:   filepath.Join(mviper.GetString("VirtualRegistryManagement.logger.system_log.path"), fileName),
		MaxSize:    mviper.GetInt("VirtualRegistryManagement.logger.system_log.max_size"),    // MB
		MaxBackups: mviper.GetInt("VirtualRegistryManagement.logger.system_log.max_backups"), // how many backup for log file.
		MaxAge:     mviper.GetInt("VirtualRegistryManagement.logger.system_log.max_age"),
		Compress:   mviper.GetBool("VirtualRegistryManagement.logger.system_log.compress"),
	}

	fileWriter := zapcore.AddSync(hook)
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	consoleCore := zapcore.NewCore(consoleEncoder, consoleDebugging, mode)
	loggerCore := zapcore.NewCore(consoleEncoder, fileWriter, mode)

	core := zapcore.NewTee(loggerCore)

	if mviper.GetBool("VirtualRegistryManagement.logger.system_log.show_in_console") {
		core = zapcore.NewTee(consoleCore, loggerCore)
	}

	logger := zap.New(core, zap.AddCaller(), zap.Development())
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

}
