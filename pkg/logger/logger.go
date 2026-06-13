package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(logPath string, level string, format string) error {
	Log = logrus.New()

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Log.SetLevel(logLevel)

	if format == "text" {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		Log.SetFormatter(&logrus.JSONFormatter{})
	}

	fileLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	Log.SetOutput(fileLogger)

	return nil
}
