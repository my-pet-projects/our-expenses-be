package config

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv_ReturnsEnvironmentVariable(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	envValue := "envValue"
	defaultValue := "default"
	os.Setenv(envKey, envValue)

	result := getEnv(envKey, defaultValue)

	if result != envValue {
		t.Errorf("Should read and return environment variable. Got: '%s', want: '%s'.", result, envValue)
	}
}

func TestGetEnv_ReturnsDefaultValue(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	defaultValue := "default"

	result := getEnv(envKey, defaultValue)

	if result != defaultValue {
		t.Errorf("Should return default value instead of environment variable. Got: '%s', want: '%s'.", result, defaultValue)
	}
}

func TestGetEnvAsInt_ReturnsEnvironmentVariable(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	envValue := 10
	defaultValue := 5
	os.Setenv(envKey, strconv.Itoa(envValue))

	result := getEnvAsInt(envKey, defaultValue)

	if result != envValue {
		t.Errorf("Should read and return environment variable as integer. Got: '%d', want: '%d'.", result, envValue)
	}
}

func TestGetEnvAsInt_ReturnsDefaultValue(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	defaultValue := 5

	result := getEnvAsInt(envKey, defaultValue)

	if result != defaultValue {
		t.Errorf("Should return integer default value instead of environment variable. Got: '%d', want: '%d'.", result, defaultValue)
	}
}

func TestGetEnvAsBool_StringBool_ReturnsEnvironmentVariable(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	envValue := true
	defaultValue := false
	os.Setenv(envKey, strconv.FormatBool(envValue))

	result := getEnvAsBool(envKey, defaultValue)

	if !result {
		t.Errorf("Should read and return environment variable as string boolean. Got: '%t', want: '%t'.", result, envValue)
	}
}

func TestGetEnvAsBool_StringInt_ReturnsEnvironmentVariable(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	envValue := "1"
	defaultValue := false
	os.Setenv(envKey, envValue)

	result := getEnvAsBool(envKey, defaultValue)

	if !result {
		t.Errorf("Should read and return environment variable as integer boolean. Got: '%t', want: '%s'.", result, envValue)
	}
}

func TestGetEnvAsBool_ReturnsDefaultValue(t *testing.T) {
	os.Clearenv()
	envKey := "envVariable"
	defaultValue := true

	result := getEnvAsBool(envKey, defaultValue)

	if result != defaultValue {
		t.Errorf("Should return boolean default value instead of environment variable. Got: '%t', want: '%t'.", result, defaultValue)
	}
}

func TestProvideConfiguration_ReturnsError(t *testing.T) {
	loadEnvError := errors.New("new error")
	// Save original function and restore it at the end.
	origLoadEnvFileFn := loadEnvFileFn
	defer func() { loadEnvFileFn = origLoadEnvFileFn }()

	// Simulate throw error.
	loadEnvFileFn = func(filenames ...string) (err error) {
		return loadEnvError
	}

	_, error := ProvideConfiguration()
	if error == nil {
		t.Error("Should throw error.")
	}

	assert.Equal(t, loadEnvError, error)
}

func TestProvideConfiguration_ReturnsConfig(t *testing.T) {
	config := &Config{
		Mongo: &MongoConfig{
			URI:      "connection_string",
			Database: "database",
		},
		Port:      80,
		LogLevel:  "INFO",
		LogToFile: true,
	}

	// Save original function and restore it at the end.
	origLoadEnvFileFn := loadEnvFileFn
	defer func() { loadEnvFileFn = origLoadEnvFileFn }()

	loadEnvFileFn = func(filenames ...string) (err error) {
		return origLoadEnvFileFn("../.env.sample")
	}

	result, error := ProvideConfiguration()
	if error != nil {
		t.Errorf("Should not throw error. Thrown error: '%s'.", error.Error())
	}

	if result == nil {
		t.Error("Should return config object.")
	}

	assert.Equal(t, result, config)
}
