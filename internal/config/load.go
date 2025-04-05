package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
)

func LoadConfig[ConfigType ServiceConfig](ctx context.Context, awsConfig aws.Config, configFilename *string) (ConfigType, error) {
	var c ConfigType
	var err error

	// try to load configuration from file
	if *configFilename != "" {
		c, err = LoadFromFile[ConfigType](&LoadFromFileInput{Filename: *configFilename})
		if err != nil {
			return c, errors.Wrap(err, "failed to load configuration from file")
		}
	}

	// try to load configuration from AppConfig
	stageName := os.Getenv("STAGE")
	appName := os.Getenv("APP_NAME")
	configProfile := os.Getenv("CONFIG_PROFILE")
	if !c.IsValid() && stageName != "" && appName != "" && configProfile != "" {
		c, err = LoadFromAppConfig[ConfigType](ctx, &LoadFromAppConfigInput{
			AWSConfig:                awsConfig,
			ApplicationName:          appName,
			ConfigurationProfileName: configProfile,
			EnvironmentName:          stageName,
		})
		if err != nil {
			return c, errors.Wrap(err, "failed to load configuration from AppConfig")
		}
	}

	if !c.IsValid() {
		return c, fmt.Errorf("no configuration file provided, nor STAGE, APP_NAME and CONFIG_PROFILE environment variables set to load from AppConfig")
	}

	return c, nil
}
