// Package config implements functionality to read and validate yaml configuration file.
package config

import (
	"io/ioutil"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// nolint:gochecknoglobals
var readFileFn = ioutil.ReadFile

// NewConfig provides application configuration based on yaml config file.
func NewConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	bytes, err := readFileFn(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read config file")
	}

	config := &Config{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, errors.Wrap(err, "decode yaml config file")
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
