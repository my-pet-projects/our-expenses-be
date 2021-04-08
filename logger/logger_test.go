package logger

import (
	"errors"
	"fmt"
	"os"
	"our-expenses-server/config"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInfo_DoNotThrowsException(t *testing.T) {
	msg := "info message"

	appLogger := &AppLogger{logrus.New()}

	appLogger.Info(msg, FieldsSet{})
	appLogger.Info(msg, FieldsSet{"key": "value"})
}

func TestError_DoNotThrowsException(t *testing.T) {
	msg := "info message"
	error := errors.New("new error")

	appLogger := &AppLogger{logrus.New()}

	appLogger.Error(msg, error, FieldsSet{})
	appLogger.Error(msg, error, FieldsSet{"key": "value"})
}

func TestFatal_DoNotThrowsException(t *testing.T) {
	msg := "info message"
	error := errors.New("new error")

	appLogger := &AppLogger{logrus.New()}

	// Define custom exit function
	defer func() { appLogger.Logger.ExitFunc = nil }()
	var fatal bool
	appLogger.Logger.ExitFunc = func(int) { fatal = true }

	appLogger.Fatal(msg, error, FieldsSet{})
	appLogger.Fatal(msg, error, FieldsSet{"key": "value"})

	assert.Equal(t, true, fatal)
}

func TestLoggerWriters_ReturnsSingleWriter_WhenLogToFileDisabled(t *testing.T) {
	config := &config.Config{LogToFile: false}

	results, error := loggerWriters(config)

	assert.Len(t, results, 1)
	assert.Contains(t, results, os.Stdout)
	assert.Nil(t, error, "Error should be nil.")
}

func TestLoggerWriters_ReturnsMultipleWriters_WhenLogToFileEnabled(t *testing.T) {
	config := &config.Config{LogToFile: true}
	file := os.NewFile(1, "file")

	// Save original function and restore it at the end.
	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()

	// Simulate throw error.
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return file, nil
	}

	results, error := loggerWriters(config)

	assert.Len(t, results, 2)
	assert.Contains(t, results, os.Stdout)
	assert.Contains(t, results, file)
	assert.Nil(t, error, "Error should be nil.")
}

func TestLoggerWriters_ThrowsError_WhenFailedToOpenFile(t *testing.T) {
	config := &config.Config{LogToFile: true}
	openFileError := errors.New("open file error")

	// Save original function and restore it at the end.
	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()

	// Simulate throw error.
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, openFileError
	}

	results, error := loggerWriters(config)

	assert.NotNil(t, error, "Error should not be nil.")
	assert.Nil(t, results, "Result should be nil.")
}

func TestGetCallerInfo_ReturnsFileNameWithLineNum(t *testing.T) {
	fileName := "file.go"
	filePath := fmt.Sprintf("/path/to/the/%s", fileName)
	lineNum := 10

	// Save original function and restore it at the end.
	origRuntimeCallerFn := runtimeCallerFn
	defer func() { runtimeCallerFn = origRuntimeCallerFn }()

	// Simulate throw error.
	runtimeCallerFn = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return 1, filePath, lineNum, true
	}

	results := getCallerInfo()

	assert.Equal(t, fmt.Sprintf("%s#%d", fileName, lineNum), results)
}

func TestGetCallerInfo_ReturnsEmptyString(t *testing.T) {
	fileName := "file.go"
	filePath := fmt.Sprintf("/path/to/the/%s", fileName)
	lineNum := 10

	// Save original function and restore it at the end.
	origRuntimeCallerFn := runtimeCallerFn
	defer func() { runtimeCallerFn = origRuntimeCallerFn }()

	// Simulate throw error.
	runtimeCallerFn = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return 1, filePath, lineNum, false
	}

	results := getCallerInfo()

	assert.Empty(t, results)
}

func TestProvideLogger_ReturnsConfig(t *testing.T) {
	config := &config.Config{LogLevel: "INFO"}

	results, error := ProvideLogger(config)

	assert.NotNil(t, results, "Result should not be nil.")
	assert.Nil(t, error, "Error should be nil.")
}

func TestProvideLogger_ThrowsError_WhenInvalidLogLevel(t *testing.T) {
	config := &config.Config{LogLevel: "invalid"}

	results, error := ProvideLogger(config)

	assert.Nil(t, results, "Result should be nil.")
	assert.NotNil(t, error, "Error should not be nil.")
}

func TestProvideLogger_ThrowsError_WhenInvalidWriters(t *testing.T) {
	config := &config.Config{LogLevel: "INFO", LogToFile: true}

	// Save original function and restore it at the end.
	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()

	// Simulate throw error.
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("open file error")
	}

	results, error := ProvideLogger(config)

	assert.Nil(t, results, "Result should be nil.")
	assert.NotNil(t, error, "Error should not be nil.")
}
