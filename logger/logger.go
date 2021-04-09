package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"dev.azure.com/filimonovga/ourexpenses/our-expenses-server/config"
	"dev.azure.com/filimonovga/ourexpenses/our-expenses-server/entity"
	"dev.azure.com/filimonovga/ourexpenses/our-expenses-server/utils"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// AppLogger is a wrapper around a Logrus logger.
type AppLogger struct {
	*logrus.Logger
}

// AppLoggerInterface defines a contract to log event.
type AppLoggerInterface interface {
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, args ...interface{})
	InfofWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{})
	InfoWithFields(ctx context.Context, msg string, fields FieldsSet)

	Error(ctx context.Context, msg string, err error)
	Errorf(ctx context.Context, format string, err error, args ...interface{})
	ErrorfWithFields(ctx context.Context, format string, err error, fields FieldsSet, args ...interface{})
	ErrorWithFields(ctx context.Context, msg string, err error, fields FieldsSet)

	AddDefaultFields(fields FieldsSet)
}

// FieldsSet is a struct that describes log entry fields.
type FieldsSet map[string]interface{}

const (
	contextLogTag string = "context"
	errorLogTag   string = "error"
	logFileName   string = "storage/logs/app.log"

	correlationIDLogField string = "correlationId"
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

	logEntry = logger.WithField("app", "our-expenses-server")

	return &AppLogger{logger}, nil
}

func (logger *AppLogger) AddDefaultFields(fields FieldsSet) {
	logEntry = logEntry.WithFields(generateLogFields(fields))
}

func (logger *AppLogger) Error(ctx context.Context, msg string, err error) {
	logger.ErrorWithFields(ctx, msg, err, nil)
}

func (logger *AppLogger) Errorf(ctx context.Context, format string, err error, args ...interface{}) {
	logger.ErrorWithFields(ctx, fmt.Sprintf(format, args...), err, nil)
}

func (logger *AppLogger) ErrorfWithFields(ctx context.Context, format string, err error, fields FieldsSet, args ...interface{}) {
	logger.ErrorWithFields(ctx, fmt.Sprintf(format, args...), err, fields)
}

// Error logs a message at level Error.
func (logger *AppLogger) ErrorWithFields(ctx context.Context, errMessage string, err error, fields FieldsSet) {

	fmt.Printf("\n ErrorWithFields:\n %+v \n\n", err)

	stackTrace := getStackTrace(err)
	cause := getCause(err)

	appError, ok := err.(entity.AppError)
	if ok {
		stackTrace = getStackTrace(appError.Err)
		cause = getCause(appError.Err)
	}

	logFields := logrus.Fields{}

	entry := logEntry.
		WithField(contextLogTag, getCallerInfo()).
		// WithField(errorLogTag, appError.Error()).
		WithFields(logFields).
		WithField("cause", cause)

	if cause == "" {
		entry = entry.WithField("cause", fmt.Sprintf("%+v", err.Error()))
	}

	if stackTrace != "" {
		entry = entry.WithField("stack", stackTrace)
	}

	entry.Error(errMessage)
}

// Info logs a message at level Info.
// func (logger *AppLogger) Info(msg string, fields map[string]interface{}) {

// 	logFields := logrus.Fields{
// 		// "application": "our-expenses-server",
// 		// "host":        os.Getenv("HOSTNAME"),
// 		// "environment":   "TBD",
// 		// "correlationId": "TBD",
// 	}

// 	logEntry = logger.WithFields(logrus.Fields{})
// 	if fields != nil {
// 		for key, val := range fields {
// 			logEntry = logEntry.WithField(key, val)
// 		}
// 	}
// 	logEntry.
// 		WithField(contextLogTag, getCallerInfo()).
// 		WithFields(logFields).
// 		Info(msg)
// }

func (logger *AppLogger) Info(ctx context.Context, msg string) {
	logger.InfoWithFields(ctx, msg, nil)
}

func (logger *AppLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	logger.InfoWithFields(ctx, fmt.Sprintf(format, args...), nil)
}

func (logger *AppLogger) InfofWithFields(ctx context.Context, format string, fields FieldsSet, args ...interface{}) {
	logger.InfoWithFields(ctx, fmt.Sprintf(format, args...), fields)
}

func (logger *AppLogger) InfoWithFields(ctx context.Context, msg string, fields FieldsSet) {
	logEntry := createLogEntry(ctx)
	logEntry.
		WithFields(generateLogFields(fields)).
		WithField("context", getCallerInfo()).
		Info(msg)
}

func createLogEntry(ctx context.Context) *logrus.Entry {
	fields := getCommonFields()

	correlationID := utils.GetContextStringValue(ctx, utils.ContextKeyCorrelationID)
	if correlationID != "" {
		fields[correlationIDLogField] = correlationID
	}

	return logrus.WithFields(fields)
}

func generateLogFields(fields FieldsSet) logrus.Fields {
	logFields := logrus.Fields{}
	for k, v := range fields {
		logFields[k] = v
	}
	return logFields
}

func getCommonFields() logrus.Fields {
	logFields := logrus.Fields{}
	// logFields[contextLogTag] = getCallerInfo()
	return logFields
}

// getCallerInfo returns an information about caller filename with line number. Example: "main.go#10".
func getCallerInfo() string {
	_, filePath, lineNo, isOk := runtimeCallerFn(3)
	if isOk {
		pathArray := strings.Split(filePath, "/")
		fileName := pathArray[len(pathArray)-1]
		return fmt.Sprintf("%s#%d", fileName, lineNo)
	}

	// notice that we're using 1, so it will actually log the where
	// the error happened, 0 = this function, we don't want that.
	// pc, fn, line, _ := runtime.Caller(1)

	// log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)

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

// ??
func getStackTrace(err error) string {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	e, ok := errors.Cause(err).(stackTracer)
	if ok {
		// stack trace could be retrieved via stackTracer interface
		st := e.StackTrace()
		return fmt.Sprintf("%+v", st)
	}

	return ""
}

func getCause(err error) string {
	type causer interface {
		Cause() error
	}

	e, ok := err.(causer)
	if ok {
		return e.Cause().Error()
	}

	return ""
	// 		return err.(causer).Cause().Error()  -works
	// return errors.Cause(err).Error()
}
