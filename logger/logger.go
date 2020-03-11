package logger

import (
	"fmt"
	"io"
	"os"
	"our-expenses-server/config"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

// AppLogger is a wrapper around a Logrus logger.
type AppLogger struct {
	*logrus.Logger
}

// Fields is a struct that describes log entry fields.
type Fields map[string]interface{}

const (
	contextLogTag string = "context"
	errorLogTag   string = "error"
	logFileName   string = "storage/logs/app.log"
)

var logEntry *logrus.Entry
var openFileFn = os.OpenFile
var runtimeCallerFn = runtime.Caller

// ProvideLogger returns a new initialized logger with the level, formatter, and output set.
func ProvideLogger(config *config.Config) (*AppLogger, error) {
	level, levelErr := logrus.ParseLevel(config.LogLevel)
	if levelErr != nil {
		return nil, levelErr
	}

	writers, writersErr := loggerWriters(config)
	if writersErr != nil {
		return nil, writersErr
	}

	logger := logrus.New()
	logger.Out = io.MultiWriter(writers...)
	logger.Level = level
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Stamp,
	})

	return &AppLogger{logger}, nil
}

// Error logs a message at level Error.
func (logger *AppLogger) Error(errMessage string, err error, fields map[string]interface{}) {
	logEntry = logger.WithFields(logrus.Fields{})
	if fields != nil {
		for key, val := range fields {
			logEntry = logEntry.WithField(key, val)
		}
	}
	logEntry.
		WithField(contextLogTag, getCallerInfo()).
		WithField(errorLogTag, err).
		Error(errMessage)
}

// Fatal logs a message at level Fatal.
func (logger *AppLogger) Fatal(errMessage string, err error, fields map[string]interface{}) {
	logEntry = logger.WithFields(logrus.Fields{})
	if fields != nil {
		for key, val := range fields {
			logEntry = logEntry.WithField(key, val)
		}
	}
	logEntry.
		WithField(contextLogTag, getCallerInfo()).
		WithField(errorLogTag, err).
		Fatal(errMessage)
}

// Info logs a message at level Info.
func (logger *AppLogger) Info(msg string, fields map[string]interface{}) {
	logEntry = logger.WithFields(logrus.Fields{})
	if fields != nil {
		for key, val := range fields {
			logEntry = logEntry.WithField(key, val)
		}
	}
	logEntry.
		WithField(contextLogTag, getCallerInfo()).
		Info(msg)
}

// getCallerInfo returns an information about caller filename with line number. Example: "main.go#10".
func getCallerInfo() string {
	_, filePath, lineNo, isOk := runtimeCallerFn(2)
	if isOk {
		pathArray := strings.Split(filePath, "/")
		fileName := pathArray[len(pathArray)-1]
		return fmt.Sprintf("%s#%d", fileName, lineNo)
	}

	return ""
}

// loggerWriters returns an array with output writes.
func loggerWriters(config *config.Config) ([]io.Writer, error) {
	writers := []io.Writer{colorable.NewColorableStdout()}

	if config.LogToFile {
		file, fileError := openFileFn(logFileName, os.O_CREATE|os.O_WRONLY, 0666)
		if fileError != nil {
			return nil, fileError
		}
		writers = append(writers, file)
	}

	return writers, nil
}
