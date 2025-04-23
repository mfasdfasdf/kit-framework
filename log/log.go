package log

import (
	"github.com/charmbracelet/log"
	"github.com/mfasdfasdf/kit-framework/config"
	"os"
	"time"
)

var logger *log.Logger

func InitLog() {
	logger = log.New(os.Stdout)
	if config.Configuration.Log.Level == "debug" {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	logger.SetPrefix(config.Configuration.AppName)
	logger.SetReportTimestamp(true)
	logger.SetTimeFormat(time.DateTime)
}

func Info(format string, values ...any) {
	if len(values) == 0 {
		logger.Infof(format)
	} else {
		logger.Infof(format, values...)
	}
}

func Warn(format string, values ...any) {
	if len(values) == 0 {
		logger.Warnf(format)
	} else {
		logger.Warnf(format, values...)
	}
}

func Error(format string, values ...any) {
	if len(values) == 0 {
		logger.Errorf(format)
	} else {
		logger.Errorf(format, values...)
	}
}

func Fatal(format string, values ...any) {
	if len(values) == 0 {
		logger.Fatalf(format)
	} else {
		logger.Fatalf(format, values...)
	}
}
