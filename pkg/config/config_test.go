package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_InvalidFilePath_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	origReadFileFn := readFileFn
	defer func() { readFileFn = origReadFileFn }()
	readFileFn = func(filename string) ([]byte, error) {
		return nil, errors.New("failed to read the file")
	}

	// Act
	cfg, cfgErr := NewConfig()

	// Assert
	assert.Nil(t, cfg, "Result should be nil.")
	assert.NotNil(t, cfgErr, "Error should not be nil.")
}

func TestNewConfig_InvalidYamlFormat_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	invalidYaml := "invalid yaml file content"

	origReadFileFn := readFileFn
	defer func() { readFileFn = origReadFileFn }()
	readFileFn = func(filename string) ([]byte, error) {
		return []byte(invalidYaml), nil
	}

	// Act
	cfg, cfgErr := NewConfig()

	// Assert
	assert.Nil(t, cfg, "Result should be nil.")
	assert.NotNil(t, cfgErr, "Error should not be nil.")
}

func TestNewConfig_ValidYamlStructure_ReturnsConfigStruct(t *testing.T) {
	t.Parallel()
	// Arrange
	validYaml := `
    application:
        name: app-name
        version: 1
    `

	origReadFileFn := readFileFn
	defer func() { readFileFn = origReadFileFn }()
	readFileFn = func(filename string) ([]byte, error) {
		return []byte(validYaml), nil
	}

	// Act
	cfg, cfgErr := NewConfig()

	// Assert
	assert.NotNil(t, cfg, "Result should not be nil.")
	assert.Nil(t, cfgErr, "Error should be nil.")
}

func TestValidate_InvalidConfig_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := new(Config)

	// Act
	validationErr := cfg.Validate()

	// Assert
	assert.NotNil(t, validationErr, "Error should not be nil.")
}

func TestValidate_ValidConfig_DoesNotThrowError(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := &Config{
		Server: Server{
			Name: "name",
			Host: "host",
			Port: 8080,
			Timeout: Timeout{
				Shutdown: 10,
				Write:    10,
				Read:     10,
				Idle:     10,
			},
			Security: Security{
				Jwt: Jwt{
					SecretKey:              "key",
					TokenExpiration:        10,
					RefreshTokenExpiration: 100,
				},
			},
		},
		Logger: Logger{
			Name:       "name",
			Level:      "level",
			JSONFormat: true,
			Writers: Writers{
				FileWriter: FileWriter{
					Enabled: true,
					Path:    "path",
				},
			},
		},
		Database: Database{
			Mongo: MongoDB{
				Name:     "name",
				URI:      "uri",
				User:     "user",
				Pass:     "pass",
				Database: "database",
			},
		},
		Telemetry: Telemetry{
			ServiceName: "name",
			Level:       "debug",
			AccessToken: "token",
		},
	}

	// Act
	validationErr := cfg.Validate()

	// Assert
	assert.Nil(t, validationErr, "Error should not be nil.")
}
