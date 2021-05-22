package logger

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger_InvalidLogLevel_ThrowsError(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Level: "invalid",
	}

	// Act
	logger, loggerErr := NewLogger(cfg)

	// Assert
	assert.Nil(t, logger, "Result should be nil.")
	assert.NotNil(t, loggerErr, "Error should not be nil.")
}

func TestNewLogger_InvalidWriters_ThrowsError(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Level: "DEBUG",
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: true,
			},
		},
	}

	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("open file error")
	}

	// Act
	result, error := NewLogger(cfg)

	// Assert
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, error, "Error should not be nil.")
}

func TestNewLogger_ReturnsLogger(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Level:      "DEBUG",
		JSONFormat: true,
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: false,
			},
		},
	}

	// Act
	result, error := NewLogger(cfg)

	// Assert
	assert.NotNil(t, result, "Result should be nil.")
	assert.Nil(t, error, "Error should not be nil.")
}

func TestLog_DoNotThrowsException(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Level: "DEBUG",
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: false,
			},
		},
	}
	ctx := context.Background()
	msg := "info message"
	msgf := "info %s message"
	f := "format"
	fields := FieldsSet{"key1": "value1", "key2": "value2"}
	err := errors.New("error")
	var buf bytes.Buffer

	// SUT
	sut, _ := NewLogger(cfg)
	sut.logger.SetOutput(&buf)

	// Act
	sut.Info(ctx, msg)
	sut.Infof(ctx, msgf, f)
	sut.InfoWithFields(ctx, msg, fields)
	sut.InfofWithFields(ctx, msgf, fields, f)
	sut.Warn(ctx, msg)
	sut.Warnf(ctx, msgf, f)
	sut.WarnWithFields(ctx, msg, fields)
	sut.WarnfWithFields(ctx, msgf, fields, f)
	sut.Error(ctx, msg, err)
	sut.Errorf(ctx, msgf, err, f)
	sut.ErrorWithFields(ctx, msg, err, fields)
	sut.ErrorfWithFields(ctx, msgf, err, fields, f)
}

func TestSetupWriters_ReturnsSingleWriter_WhenLogToFileDisabled(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: false,
			},
		},
	}

	// Act
	results, error := setupWriters(cfg.Writers)

	// Assert
	assert.Len(t, results, 1)
	assert.Contains(t, results, os.Stdout)
	assert.Nil(t, error, "Error should be nil.")
}

func TestSetupWriters_ReturnsMultipleWriters_WhenLogToFileEnabled(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: true,
			},
		},
	}
	file := os.NewFile(1, "file")

	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return file, nil
	}

	// Act
	results, error := setupWriters(cfg.Writers)

	// Assert
	assert.Len(t, results, 2)
	assert.Contains(t, results, os.Stdout)
	assert.Contains(t, results, file)
	assert.Nil(t, error, "Error should be nil.")
}

func TestLoggerWriters_ThrowsError_WhenFailedToOpenFile(t *testing.T) {
	// Arrange
	cfg := config.Logger{
		Writers: config.Writers{
			FileWriter: config.FileWriter{
				Enabled: true,
			},
		},
	}

	origOpenFileFn := openFileFn
	defer func() { openFileFn = origOpenFileFn }()
	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("open file error")
	}

	// Act
	results, error := setupWriters(cfg.Writers)

	// Assert
	assert.Nil(t, results, "Result should be nil.")
	assert.NotNil(t, error, "Error should not be nil.")
}

func TestGetCallerInfo_CanGetCallerInfo_ReturnsFileNameWithLineNum(t *testing.T) {
	// Arrange
	fileName := "file.go"
	lineNum := 10

	origRuntimeCallerFn := runtimeCallerFn
	defer func() { runtimeCallerFn = origRuntimeCallerFn }()
	runtimeCallerFn = func(skip int) (pc uintptr, file string, line int, ok bool) {
		filePath := fmt.Sprintf("/path/to/the/%s", fileName)
		return 1, filePath, lineNum, true
	}

	// Act
	results := getCallerInfo()

	// Assert
	assert.Equal(t, fmt.Sprintf("%s#%d", fileName, lineNum), results)
}

func TestGetCallerInfo_FailedGetCallerInfo_ReturnsFileNameWithLineNum(t *testing.T) {
	// Arrange
	fileName := "file.go"
	lineNum := 10

	origRuntimeCallerFn := runtimeCallerFn
	defer func() { runtimeCallerFn = origRuntimeCallerFn }()
	runtimeCallerFn = func(skip int) (pc uintptr, file string, line int, ok bool) {
		filePath := fmt.Sprintf("/path/to/the/%s", fileName)
		return 1, filePath, lineNum, false
	}

	// Act
	results := getCallerInfo()

	// Assert
	assert.Empty(t, results)
}
