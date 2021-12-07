// Package config implements functionality to read and validate yaml configuration file.
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// nolint:gochecknoglobals
var (
	readFileFn = ioutil.ReadFile
	getwdFn    = os.Getwd
)

const configPath = "config/config.yaml"

// NewConfig provides application configuration based on yaml config file.
func NewConfig() (*Config, error) {
	workingDir, _ := getwdFn()
	bytes, bytesErr := readFileFn(filepath.Join(workingDir, configPath))
	if bytesErr != nil {
		return nil, errors.Wrap(bytesErr, "read config file")
	}

	config := &Config{}
	if unmarshallErr := yaml.Unmarshal(bytes, config); unmarshallErr != nil {
		return nil, errors.Wrap(unmarshallErr, "decode yaml config file")
	}

	return config, nil
}

// Validate validates config struct.
func (cfg Config) Validate() error {
	validator := validator.New()
	if validatorErr := validator.Struct(cfg); validatorErr != nil {
		return errors.Wrap(validatorErr, "invalid config")
	}

	return nil
}
