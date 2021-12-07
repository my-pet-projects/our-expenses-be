package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/apperror"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
)

// Logger is a wrapper around Logrus logger.
type Logger struct {
	logger *logrus.Logger
}

// LogInterface defines a contract to log event.
type LogInterface interface {
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, args ...interface{})
	InfofWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{})
	InfoWithFields(ctx context.Context, msg string, fields FieldsSet)

	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, args ...interface{})
	WarnfWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{})
	WarnWithFields(ctx context.Context, msg string, fields FieldsSet)

	Error(ctx context.Context, msg string, err error)
	Errorf(ctx context.Context, format string, err error, args ...interface{})
	ErrorfWithFields(ctx context.Context, format string, err error, fields FieldsSet, args ...interface{})
	ErrorWithFields(ctx context.Context, msg string, err error, fields FieldsSet)
}

// FieldsSet is a struct that describes log entry fields.
type FieldsSet map[string]interface{}

// nolint:gochecknoglobals
var (
	logEntry        *logrus.Entry
	openFileFn      = os.OpenFile
	runtimeCallerFn = runtime.Caller
)

// NewLogger returns a new initialized logger with the level, formatter, and output set.
func NewLogger(config config.Logger) (*Logger, error) {
	level, levelErr := logrus.ParseLevel(config.Level)
	if levelErr != nil {
		return nil, errors.Wrap(levelErr, "level parse")
	}

	writers, writersErr := setupWriters(config.Writers)
	if writersErr != nil {
		return nil, writersErr
	}

	logger := logrus.New()
	logger.Out = io.MultiWriter(writers...)
	logger.Level = level

	if config.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: "severity",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.Stamp,
		})
	}

	logEntry = logger.
		// WithField("host", os.Getenv("HOSTNAME")).
		// WithField("environment", "TBD").
		WithField("app", config.Name)

	return &Logger{logger}, nil
}

// Info logs a message at level Info.
func (logger *Logger) Info(ctx context.Context, msg string) {
	log(ctx, logrus.InfoLevel, msg, nil, nil)
}

// Infof logs a formatted message at level Info.
func (logger *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	log(ctx, logrus.InfoLevel, fmt.Sprintf(format, args...), nil, nil)
}

// InfofWithFields logs a formatted message with additional fields at level Info.
func (logger *Logger) InfofWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{}) {
	log(ctx, logrus.InfoLevel, fmt.Sprintf(format, args...), nil, fields)
}

// InfoWithFields logs a message with additional fields at level Info.
func (logger *Logger) InfoWithFields(ctx context.Context, msg string, fields FieldsSet) {
	log(ctx, logrus.InfoLevel, msg, nil, fields)
}

// Warn logs a message at level Info.
func (logger *Logger) Warn(ctx context.Context, msg string) {
	log(ctx, logrus.WarnLevel, msg, nil, nil)
}

// Warnf logs a formatted message at level Info.
func (logger *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	log(ctx, logrus.WarnLevel, fmt.Sprintf(format, args...), nil, nil)
}

// WarnfWithFields logs a formatted message with additional fields at level Info.
func (logger *Logger) WarnfWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{}) {
	log(ctx, logrus.WarnLevel, fmt.Sprintf(format, args...), nil, fields)
}

// WarnWithFields logs a message with additional fields at level Info.
func (logger *Logger) WarnWithFields(ctx context.Context, msg string, fields FieldsSet) {
	log(ctx, logrus.WarnLevel, msg, nil, fields)
}

// Error logs a message at level Error.
func (logger *Logger) Error(ctx context.Context, msg string, err error) {
	log(ctx, logrus.ErrorLevel, msg, err, nil)
}

// Errorf logs a formatted message at level Error.
func (logger *Logger) Errorf(ctx context.Context, format string, err error, args ...interface{}) {
	log(ctx, logrus.ErrorLevel, fmt.Sprintf(format, args...), err, nil)
}

// ErrorfWithFields logs a formatted message with additional fields at level Error.
func (logger *Logger) ErrorfWithFields(ctx context.Context, format string, err error,
	fields FieldsSet, args ...interface{}) {
	log(ctx, logrus.ErrorLevel, fmt.Sprintf(format, args...), nil, fields)
}

// ErrorWithFields logs a message with additional fields at level Error.
func (logger *Logger) ErrorWithFields(ctx context.Context, msg string, err error, fields FieldsSet) {
	log(ctx, logrus.ErrorLevel, msg, err, fields)
}

func log(ctx context.Context, level logrus.Level, msg string, err error, fields FieldsSet) {
	logEntry.
		WithFields(buildFieldsSetLogFields(fields)).
		WithFields(buildCommonLogFields(ctx)).
		WithFields(buildErrorLogFields(err)).
		Log(level, msg)
}

// buildCommonLogFields returns common log fields such as correlationId, context and etc.
func buildCommonLogFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	fields["correlationId"] = "TBD"
	fields["context"] = getCallerInfo()
	return fields
}

// buildFieldsSetLogFields returns log fields from FieldsSet.
func buildFieldsSetLogFields(fieldsSet FieldsSet) logrus.Fields {
	fields := logrus.Fields{}
	for k, v := range fieldsSet {
		fields[k] = v
	}
	return fields
}

// buildErrorLogFields returns error specific log fields.
func buildErrorLogFields(err error) logrus.Fields {
	if err == nil {
		return logrus.Fields{}
	}
	fields := logrus.Fields{}
	fields["cause"] = apperror.GetCause(err)
	fields["stack"] = apperror.GetStackTrace(err)
	return fields
}

// getCallerInfo returns information about caller filename with line number. Example: "main.go#10".
func getCallerInfo() string {
	// We need to skip 4 callers:
	// 0 = this function
	// 1 = generateLogFields function
	// 2 = Info*/Error* functions
	// 3 = log function
	// 4 = original caller
	skip := 4
	_, filePath, lineNo, isOk := runtimeCallerFn(skip)
	if isOk {
		pathArray := strings.Split(filePath, "/")
		fileName := pathArray[len(pathArray)-1]
		return fmt.Sprintf("%s#%d", fileName, lineNo)
	}

	return ""
}

// setupWriters returns an array with output writes.
func setupWriters(config config.Writers) ([]io.Writer, error) {
	writers := []io.Writer{
		colorable.NewColorableStdout(),
	}

	if config.FileWriter.Enabled {
		file, fileError := openFileFn(config.FileWriter.Path, os.O_CREATE|os.O_WRONLY, 0666)
		if fileError != nil {
			return nil, fileError
		}
		writers = append(writers, file)
	}

	return writers, nil
}
