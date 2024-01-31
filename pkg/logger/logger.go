package logger

import (
	"fmt"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type MyLogger struct {
	*logrus.Logger
}

func (l *MyLogger) Error(message string, err ...string) {
	fmt.Printf("%s: %s", message, err)
}

func Error(err error) string {
	return err.Error()
}

func Any(key string, val any) string {
	return fmt.Sprintf("%s: %v", key, val)
}

func String(key string, val string) string {
	return fmt.Sprintf("%s: %v", key, val)
}

func NewLogger() *MyLogger {
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "log/info.log",
		logrus.ErrorLevel: "log/error.log",
	}

	Log := logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return &MyLogger{
		Logger: Log,
	}
}
