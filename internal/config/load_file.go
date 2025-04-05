package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type LoadFromFileInput struct {
	Filename string
}

func LoadFromFile[ConfigType ServiceConfig](input *LoadFromFileInput) (ConfigType, error) {
	//#nosec: G304 (CWE-22): Potential file inclusion via variable
	var config ConfigType
	configFile, err := os.Open(input.Filename)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %s", err.Error())
	}

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %s", err.Error())
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file: %s", err.Error())
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("invalid configuration: %s", err.Error())
	}

	return config, nil
}
