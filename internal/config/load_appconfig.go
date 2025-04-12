package config

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type LoadFromAppConfigInput struct {
	AWSConfig                aws.Config
	ApplicationName          string
	ConfigurationProfileName string
	EnvironmentName          string
	DefaultConfig            ServiceConfig
}

// LoadFromAppConfig loads the jarbas configuration from AWS AppConfig.
func LoadFromAppConfig[ConfigType ServiceConfig](ctx context.Context, input *LoadFromAppConfigInput) (ConfigType, error) {
	log.Ctx(ctx).Debug().Interface("input", input).Msg("loading configuration from AppConfig")
	appConfigClient := appconfigdata.NewFromConfig(input.AWSConfig)

	var config ConfigType
	if input.DefaultConfig != nil {
		config = input.DefaultConfig.(ConfigType)
	}

	configSessionOutput, err := appConfigClient.StartConfigurationSession(ctx, &appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:          aws.String(input.ApplicationName),
		EnvironmentIdentifier:          aws.String(input.EnvironmentName),
		ConfigurationProfileIdentifier: aws.String(input.ConfigurationProfileName),
	})
	if err != nil {
		return config, errors.Wrap(err, "failed to start config session on AppConfig")
	}

	getConfigOutput, err := appConfigClient.GetLatestConfiguration(ctx, &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: configSessionOutput.InitialConfigurationToken,
	})
	if err != nil {
		return config, errors.Wrap(err, "failed to get latest config from AppConfig")
	}

	err = json.Unmarshal(getConfigOutput.Configuration, &config)
	if err != nil {
		return config, errors.Wrap(err, "failed to unmarshal AppConfig configuration")
	}

	return config, nil
}
