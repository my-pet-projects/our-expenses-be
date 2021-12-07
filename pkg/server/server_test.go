package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewServer_ReturnsServerInstanceWithSettingsFromConfig(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Server{
		Host: "localhost",
		Port: 1234,
		Timeout: config.Timeout{
			Shutdown: 1,
			Write:    2,
			Read:     3,
			Idle:     4,
		},
	}
	logger := new(mocks.LogInterface)
	handlerFn := func(e *echo.Echo) {}

	// Act
	result := NewServer(logger, config, handlerFn)

	// Assert
	assert.NotNil(t, result, "Result should be nil.")
	assert.Equal(t, fmt.Sprintf("%s:%d", config.Host, config.Port), result.httpServer.Server.Addr, "Server address should have port from config.")
	assert.Equal(t, time.Duration(config.Timeout.Read)*time.Second, result.httpServer.Server.ReadTimeout, "Read timeout should have value from config.")
	assert.Equal(t, time.Duration(config.Timeout.Write)*time.Second, result.httpServer.Server.WriteTimeout, "Write timeout should have value from config.")
	assert.Equal(t, time.Duration(config.Timeout.Idle)*time.Second, result.httpServer.Server.IdleTimeout, "Idle timeout should have value from config.")
}

func TestStart_StartsAndGracefullyStopsServer_ShouldNotThrowError(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := config.Config{
		Server: config.Server{
			Port: 1234,
			Timeout: config.Timeout{
				Shutdown: 1,
				Write:    2,
				Read:     3,
				Idle:     4,
			},
		},
		Logger: config.Logger{
			Level:   "DEBUG",
			Writers: config.Writers{FileWriter: config.FileWriter{Enabled: false}},
		},
	}

	logger := new(mocks.LogInterface)
	ctx, cancel := context.WithCancel(context.Background())
	handlerFn := func(e *echo.Echo) {}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()

	// SUT
	sut := NewServer(logger, cfg.Server, handlerFn)

	// Act
	runningChan := make(chan struct{})
	doneChan := make(chan struct{})
	errChan := make(chan error)

	go func() {
		close(runningChan)
		errChan <- sut.Start(ctx)
		defer close(doneChan)
	}()

	<-runningChan
	cancel()
	err := <-errChan
	<-doneChan

	// Assert
	assert.Nil(t, err, "Error should be nil.")
}

func TestStart_InvalidConfig_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := config.Config{
		Server: config.Server{
			Port: -1234,
			Timeout: config.Timeout{
				Shutdown: 1,
				Write:    2,
				Read:     3,
				Idle:     4,
			},
		},
	}

	logger := new(mocks.LogInterface)
	ctx := context.Background()
	handlerFn := func(e *echo.Echo) {}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()

	// SUT
	sut := NewServer(logger, cfg.Server, handlerFn)

	// Act
	err := sut.Start(ctx)

	// Assert
	assert.NotNil(t, err, "Error should not be nil.")
}
