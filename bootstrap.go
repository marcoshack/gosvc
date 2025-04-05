package gosvc

import (
	"context"

	"github.com/marcoshack/gosvc/internal/config"
)

type Bootstrap[ConfigType config.Validatable] struct {
	Name           string
	Config         ConfigType
	ConfigFileName string
}

type BootstrapInput struct {
	Name           string
	ConfigFileName string
}

func NewBootstrap[ConfigType config.Validatable](ctx context.Context, input BootstrapInput) (*Bootstrap[ConfigType], error) {

	config, err := config.LoadFromFile[ConfigType](&config.LoadFromFileInput{
		Filename: input.ConfigFileName,
	})
	if err != nil {
		return nil, err
	}

	return &Bootstrap[ConfigType]{
		Name:           input.Name,
		ConfigFileName: input.ConfigFileName,
		Config:         config,
	}, nil
}
