package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

type LoadFromFileInput struct {
	FileName      string
	DefaultConfig ServiceConfig
}

func LoadFromFile[ConfigType ServiceConfig](input *LoadFromFileInput) (ConfigType, error) {
	//#nosec: G304 (CWE-22): Potential file inclusion via variable
	var config ConfigType
	if input.DefaultConfig != nil {
		config = input.DefaultConfig.(ConfigType)
	}

	configFile, err := os.Open(input.FileName)
	if err != nil {
		return config, errors.Wrap(err, "error opening config file")
	}

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return config, errors.Wrap(err, "error reading config file")
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return config, errors.Wrap(err, "error parsing config file")
	}

	if err := config.Validate(); err != nil {
		return config, errors.Wrap(err, "error validating config file")
	}

	return config, nil
}
